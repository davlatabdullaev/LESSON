package handler

import (
	"bazaar/api/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EndSell godoc
// @Router           /end_sell/{id} [PUT]
// @Summary          end sell
// @Description      end sell
// @Tags             sell
// @Accept           json
// @Produce          json
// @Param            id path string true "sale_id"
// @Param            SaleRequest body models.SaleRequest true "sale request"
// @Succes           200 {object} models.Response
// @Failure          400 {object} models.Response
// @Failure          404 {object} models.Response
// @Failure          500 {object} models.Response
func (h Handler) EndSale(c *gin.Context) {

	var (
		totalPrice float64
		err        error
	)

	id := c.Param("id")

	request := models.SaleRequest{}

	if err = c.ShouldBindJSON(&request); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	baskets, err := h.storage.Basket().GetList(context.Background(), models.GetBasketsListRequest{
		Page:   1,
		Limit:  100,
		Search: id,
	})

	if err != nil {
		handleResponse(c, "error while getting baskets list", http.StatusInternalServerError, err.Error())
		return
	}

	selectedProducts := make(map[string]models.Basket)

	for _, basket := range baskets.Baskets {
		totalPrice += basket.Price
		selectedProducts[basket.ProductID] = basket
		
	}
	fmt.Println(totalPrice)

	if request.Status == "cancel" {
		totalPrice = 0
		return
	}

	saleID, err := h.storage.Sale().UpdateSalePrice(context.Background(), models.SaleRequest{
		ID:         id,
		TotalPrice: int(totalPrice),
		Status:     request.Status,
	})
	if err != nil {
		handleResponse(c, "error while updating sale price and status by id", http.StatusInternalServerError, err.Error())
		return
	}

	sale, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: saleID,
	})
	if err != nil {
		handleResponse(c, "error while get sale by id", http.StatusInternalServerError, err.Error())
		return
	}

	storageData, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 100,
	})
	if err != nil {
		handleResponse(c, "error while getting storages list", http.StatusInternalServerError, err.Error())
		return
	}

	storageMap := make(map[string]models.Storage)
	for _, storage := range storageData.Storages {
		storageMap[storage.ID] = storage
	}

	for i, value := range storageMap {
		if value.ProductID == selectedProducts[value.ProductID].ProductID {
			_, err := h.storage.Storage().Update(context.Background(), models.UpdateStorage{
				ID:        i,
				ProductID: value.ProductID,
				BranchID:  value.BranchID,
				Count:     value.Count - selectedProducts[value.ProductID].Quantity,
			})

			if err != nil {
				handleResponse(c, "error while updating repositoryData prod quantities", http.StatusInternalServerError, err.Error())
				return
			}

			_, err = h.storage.StorageTransaction().Create(context.Background(), models.CreateStorageTransaction{
				ProductID:              value.ProductID,
				StaffID:                sale.CashierID,
				StorageTransactionType: "minus",
				Price:                  selectedProducts[value.ProductID].Price,
				Quantity:               float64(selectedProducts[value.ProductID].Quantity),
			})
			if err != nil {
				handleResponse(c, "error while creating storage data", http.StatusInternalServerError, err.Error())
				return
			}

		}
	}

	salesResponse, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: id,
	})

	if err != nil {
		handleResponse(c, "error while getting sales list", http.StatusInternalServerError, err.Error())
		return
	}

	if salesResponse.Status == "succes" {

		cashierResponse, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
			ID: sale.CashierID,
		})

		if err != nil {
			handleResponse(c, "error while get staff data", http.StatusInternalServerError, err.Error())
			return
		}

		cashierTarifResponse, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
			ID: cashierResponse.TarifID,
		})
		if err != nil {
			handleResponse(c, "error while getting tarif by id", http.StatusInternalServerError, err.Error())
			return
		}

		shopAssistantResponse, err := h.storage.Staff().Get(context.Background(), models.PrimaryKey{
			ID: sale.CashierID,
		})

		if err != nil {
			handleResponse(c, "error while get staff data", http.StatusInternalServerError, err.Error())
			return
		}

		shopAssistantTarifResponse, err := h.storage.Tarif().Get(context.Background(), models.PrimaryKey{
			ID: shopAssistantResponse.TarifID,
		})
		if err != nil {
			handleResponse(c, "error while getting tarif by id", http.StatusInternalServerError, err.Error())
			return
		}

		amount := 0.0

		if cashierTarifResponse.TarifType == "fixed" {

			if salesResponse.PaymentType == "card" {
				amount = (cashierTarifResponse.AmountForCard)
			} else {
				amount = (cashierTarifResponse.AmountForCash)
			}

		} else {

			if salesResponse.
				PaymentType == "card" {
				amount = (cashierTarifResponse.AmountForCard) * totalPrice
			} else {
				amount = (cashierTarifResponse.AmountForCash) * totalPrice
			}

		}

		reqToUpdate := models.UpdateStaffBalanceAndCreateTransaction{
			UpdateCashierBalance: models.StaffInfo{
				StaffID: cashierResponse.ID,
				Amount:  amount,
			},
			SaleID:          id,
			TransactionType: "topup",
			SourceType:      "sales",
			Amount:          salesResponse.Price,
			Description:     "qwerty",
		}

		if salesResponse.ShopAssistantID != "" {

			if shopAssistantTarifResponse.TarifType == "fixed" {

				if salesResponse.PaymentType == "card" {
					amount = (shopAssistantTarifResponse.AmountForCard)
				} else {
					amount = (shopAssistantTarifResponse.AmountForCash)
				}

			} else {

				if salesResponse.PaymentType == "card" {
					amount = (shopAssistantTarifResponse.AmountForCard) * totalPrice
				} else {
					amount = (shopAssistantTarifResponse.AmountForCash) * totalPrice

				}

			}

			reqToUpdate.UpdateShopAssistantBalance.StaffID = shopAssistantResponse.ID
			
			reqToUpdate.UpdateShopAssistantBalance.Amount = amount
		}

		err = h.storage.Transaction().UpdateStaffBalanceAndCreateTransaction(context.Background(), reqToUpdate)
		if err != nil {
			handleResponse(c, "error while update cashoier balance", http.StatusInternalServerError, err.Error())
			return
		}

	}
}

// barcode

package handler

import (
	"context"
	"net/http"
	"bazaar/api/models"

	"github.com/gin-gonic/gin"
)

// Barcode godoc
// @Router       /barcode [POST]
// @Summary      barcode
// @Description  barcode
// @Tags         barcode
// @Accept       json
// @Produce      json
// @Param		 info body models.Barcode true "info"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) Barcode(c *gin.Context) {
	info := models.Barcode{}
	if err := c.ShouldBindJSON(&info); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	sale, err := h.storage.Sale().Get(context.Background(), models.PrimaryKey{
		ID: info.SaleID,
	})
	if err != nil {
		handleResponse(c, "error is getting sale by id", http.StatusInternalServerError, err.Error())
		return
	}

	if sale.Status == "success" {
		handleResponse(c, "sale ended", 300, "sale ended cannot add product")
		return
	}

	if sale.Status == "cancel" {
		handleResponse(c, "sale canceled", 300, "sale canceled cannot add product")
		return
	}

	products, err := h.storage.Product().GetList(context.Background(), models.ProductGetListRequest{
		Page:    1,
		Limit:   10,
		Barcode: info.Barcode,
	})

	if err != nil {
		handleResponse(c, "error is while getting product list by barcode", http.StatusInternalServerError, err.Error())
		return
	}

	var (
		prodID    string
		prodPrice int
	)
	for _, product := range products.Products {
		prodID = product.ID
		prodPrice = int(product.Price)
	}

	baskets, err := h.storage.Basket().GetList(context.Background(), models.GetBasketsListRequest{
		Page:   1,
		Limit:  10,
		Search: info.SaleID,
	})
	if err != nil {
		handleResponse(c, "error is while getting basket list", http.StatusInternalServerError, err.Error())
		return
	}

	var (
		basketsMap = make(map[string]models.Basket)
		totalPrice = 0
	)

	totalPrice = info.Count * prodPrice

	for _, basket := range baskets.Baskets {
		basketsMap[basket.ProductID] = basket
	}

	storage, err := h.storage.Storage().GetList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 100,
	})
	if err != nil {
		handleResponse(c, "error is while getting repo list", http.StatusInternalServerError, err.Error())
		return
	}

	for _, r := range storage.Storages {
		if prodID == basketsMap[r.ProductID].ProductID {
			if r.Count < (basketsMap[r.ProductID].Quantity + info.Count) {
				handleResponse(c, "not enough product", 301, "not enough product")
				return
			}
		}

		if r.Count < info.Count {
			handleResponse(c, "not enough product", 300, "not enough product")
			return
		}
	}

	isTrue := false

	for _, value := range basketsMap {
		if prodID == value.ProductID {
			isTrue = true
			id, err := h.storage.Basket().Update(context.Background(), models.UpdateBasket{
				ID:        value.ID,
				SaleID:    value.SaleID,
				ProductID: prodID,
				Quantity:  value.Quantity + info.Count,
				Price:     value.Price + float64(totalPrice),
			})
			if err != nil {
				handleResponse(c, "error is while updating basket", 500, err.Error())
				return
			}
			updatedBasket, err := h.storage.Basket().Get(context.Background(), models.PrimaryKey{ID: id})
			if err != nil {
				handleResponse(c, "error is while getting basket", 500, err.Error())
				return
			}
			handleResponse(c, "updated", http.StatusOK, updatedBasket)
		}
	}

	if !isTrue {
		id, err := h.storage.Basket().Create(context.Background(), models.CreateBasket{
			SaleID:    info.SaleID,
			ProductID: prodID,
			Quantity:  info.Count,
			Price:     float64(totalPrice),
		})
		if err != nil {
			handleResponse(c, "error is while creating basket", 500, err.Error())
			return
		}
		createdBasket, err := h.storage.Basket().Get(context.Background(), models.PrimaryKey{ID: id})
		if err != nil {
			handleResponse(c, "error is while getting basket", 500, err.Error())
			return
		}
		handleResponse(c, "updated", http.StatusOK, createdBasket)
	}
}


/*



sale_id, barcode, count


barcode orqali prod get

count*price=batot

ba = cou

sal = said



*/
