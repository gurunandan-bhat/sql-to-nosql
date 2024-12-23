package reldb

import (
	"time"
)

type OrderSummary struct {
	IOrdID         uint32    `db:"iOrdID" json:"iOrdID,omitempty"`
	DtDt           time.Time `db:"dtDt" json:"dtDt,omitempty"`
	VRecpName      *string   `db:"vRecpName" json:"vRecpName,omitempty"`
	VRecpCountry   *string   `db:"vRecpCountry" json:"vRecpCountry,omitempty"`
	FTotal         float32   `db:"fTotal" json:"fTotal,omitempty"`
	CPaymentStatus *string   `db:"cPaymentStatus" json:"cPaymentStatus,omitempty"`
	CStatus        *string   `db:"cStatus" json:"cStatus,omitempty"`
	VSource        *string   `db:"vSource" json:"vSource,omitempty"`
}

type OrderRecipient struct {
	VRecpName        *string `db:"vRecpName" json:"vRecpName,omitempty" schema:"vRecpName,required"`
	VRecpEmail       *string `db:"vRecpEmail" json:"vRecpEmail,omitempty" schema:"-"`
	VRecpAddress     *string `db:"vRecpAddress" json:"vRecpAddress,omitempty" schema:"vRecpAddress,required"`
	VRecpAddlAddress *string `db:"-" json:"-" schema:"vRecpAddlAddress"`
	VRecpCity        *string `db:"vRecpCity" json:"vRecpCity,omitempty" schema:"vRecpCity,required"`
	VRecpState       *string `db:"vRecpState" json:"vRecpState,omitempty" schema:"vRecpState,required"`
	VRecpCountryName *string `db:"vRecpCountryName,omitempty" json:"vRecpCountryName,omitempty" schema:"vRecpCountryName,required"`
	VRecpPincode     *string `db:"vRecpPincode,omitempty" json:"vRecpPincode,omitempty" schema:"vRecpPincode,required"`
	VRecpPhone       *string `db:"vRecpPhone,omitempty" json:"vRecpPhone,omitempty" schema:"vRecpPhone,required"`
	VRecpLandLine    *string `db:"vRecpLandLine,omitempty" json:"vRecpLandLine,omitempty" schema:"-"`
}

type OrderBilling struct {
	VBillName        *string `db:"vBillName" json:"vBillName,omitempty"`
	VBillEmail       *string `db:"vBillEmail" json:"vBillEmail,omitempty"`
	VBillAddress     *string `db:"vBillAddress" json:"vBillAddress,omitempty"`
	VBillCity        *string `db:"vBillCity" json:"vBillCity,omitempty"`
	VBillState       *string `db:"vBillState" json:"vBillState,omitempty"`
	VBillCountryName *string `db:"vBillCountryName,omitempty" json:"vBillCountryName,omitempty"`
	VBillPincode     *string `db:"vBillPincode,omitempty" json:"vBillPincode,omitempty"`
	VBillPhone       *string `db:"vBillPhone,omitempty" json:"vBillPhone,omitempty"`
	VBillLandLine    *string `db:"vBillLandLine,omitempty" json:"vBillLandLine,omitempty"`
}

type OrderPayment struct {
	FTotalRate       float32 `db:"fTotalRate" json:"fTotalRate,omitempty"`
	FTotalDiscount   float32 `db:"fTotalDiscount" json:"fTotalDiscount,omitempty"`
	FShipping        float32 `db:"fShipping" json:"fShipping,omitempty"`
	FTotal           float32 `db:"fTotal" json:"fTotal,omitempty"`
	VInstructions    *string `db:"vInstructions" json:"vInstructions,omitempty"`
	VRemarks         *string `db:"vRemarks" json:"vRemarks,omitempty"`
	CPaymentType     *string `db:"cPaymentType" json:"cPaymentType,omitempty"`
	VBankName        *string `db:"vBankName" json:"vBankName,omitempty"`
	VChequeNo        *string `db:"vChequeNo" json:"vChequeNo,omitempty"`
	IInvoiceCode     uint    `db:"iInvoiceCode,omitempty" json:"iInvoiceCode,omitempty"`
	CPaymentStatus   *string `db:"cPaymentStatus" json:"cPaymentStatus,omitempty"`
	CStatus          *string `db:"cStatus" json:"cStatus,omitempty"`
	IDiscCoupID      int     `db:"iDiscCoupID" json:"iDiscCoupID,omitempty"`
	FCouponDisc      *string `db:"fCouponDisc,omitempty" json:"fCouponDisc,omitempty"`
	FGiftWrapCharges float32 `db:"fGiftWrapCharges" json:"fGiftWrapCharges,omitempty"`
}

type OrderGiftOptions struct {
	IGWCardID        uint32  `db:"iGWCardID" json:"iGWCardID,omitempty"`
	VGiftCardName    *string `db:"vGiftCardName" json:"vGiftCardName,omitempty"`
	VGiftWrapMessage *string `db:"vGiftWrapMessage" json:"vGiftWrapMessage,omitempty"`
}
type Order struct {
	IOrdID uint32    `db:"iOrdID" json:"iOrdID,omitempty"`
	DtDt   time.Time `db:"dtDt,omitempty" json:"dtDt,omitempty"`
	OrderRecipient
	OrderBilling
	OrderPayment
	OrderGiftOptions
	OrderShipments []OrderShipment `db:"-" json:"orderShipments"`
	OrderProducts  []OrderProduct  `db:"-" json:"orderProducts"`
}

type OrderProduct struct {
	IProdID         uint32  `db:"iProdID" json:"iProdID,omitempty"`
	VURLName        *string `db:"vUrlName" json:"vUrlName,omitempty"`
	VName           string  `db:"vName" json:"vName,omitempty"`
	VCategoryName   string  `db:"vCategoryName" json:"vCategoryName,omitempty"`
	VAttributeValue *string `db:"vAttributeValue" json:"vAttributeValue,omitempty"`
	VColorValue     *string `db:"vColorValue" json:"vColorValue,omitempty"`
	IQty            int     `db:"iQty" json:"iQty,omitempty"`
	FRate           float32 `db:"fRate" json:"fRate,omitempty"`
	FPayable        float32 `db:"fPayable" json:"fPayable,omitempty"`
	CGiftWrap       *string `db:"cGiftWrap" json:"cGiftWrap,omitempty"`
}

type OrderShipment struct {
	IOrdShipID   uint32    `db:"iOrdShipID" json:"iOrdShipID,omitempty"`
	IOrdID       uint32    `db:"iOrdID" json:"iOrdID,omitempty"`
	ICourierID   uint32    `db:"iCourierID" json:"iCourierID,omitempty"`
	VCourierName *string   `db:"vCourierName" json:"vCourierName,omitempty"`
	VLink        *string   `db:"vLink" json:"vLink,omitempty"`
	VShipCode    *string   `db:"vShipCode" json:"vShipCode,omitempty"`
	DtUpdatedAt  time.Time `db:"dtUpdatedAt" json:"dtUpdatedAt,omitempty"`
}

func (m *Model) OrderCount() (int, error) {

	var count int

	query := "SELECT COUNT(*) FROM orders"
	if err := m.Get(&count, query); err != nil {
		return 0, err
	}

	return count, nil
}

func (m *Model) Orders(offset, limit uint) ([]OrderSummary, error) {

	query := `SELECT
				iOrdID,
				dtDt,
				vRecpName,
				p.vCountryName as vRecpCountry,
				fTotal, 
				cPaymentStatus, 
				cStatus,
				vSource
			FROM orders
				JOIN postage p ON iRecpCountyID = p.iPostID
			ORDER BY dtDt desc
			LIMIT ?, ?`

	orders := []OrderSummary{}
	if err := m.Select(&orders, query, offset, limit); err != nil {
		return nil, err
	}

	return orders, nil
}

func (m *Model) OrderDetail(iOrdID uint) (Order, error) {

	query := `SELECT
				iOrdID,
				dtDt,
				vRecpName,
				vRecpAddress,
				vRecpCity,
				vRecpState,
				p1.vCountryName AS vRecpCountryName,
				vRecpPincode,
				vRecpPhone,
				vRecpLandLine,
				vRecpEmail,
				vBillName,
				vBillAddress,
				vBillCity,
				vBillState,
				p2.vCountryName AS vBillCountryName,
				vBillPincode,
				vBillPhone,
				vBillLandLine,
				vBillEmail,
				fTotalRate,
				fTotalDiscount,
				fShipping,
				fTotal,
				vInstructions,
				vRemarks,
				cPaymentType,
				vBankName,
				vChequeNo,
				iInvoiceCode,
				cPaymentStatus,
				o.cStatus,
				iDiscCoupID,
				fCouponDisc,
				fGiftWrapCharges,
				vGiftWrapMessage,
				o.iGWCardID,
				gwc.vName AS vGiftCardName
			FROM orders o
				LEFT JOIN postage p1 ON o.iRecpCountyID = p1.iPostID
				LEFT JOIN postage p2 ON o.iBillCountyID = p2.iPostID
				LEFT JOIN giftwrap_cards gwc ON o.iGWCardID = gwc.iGWCardID
			WHERE iOrdID = ?`

	row := m.QueryRowx(query, iOrdID)
	order := Order{}
	if err := row.StructScan(&order); err != nil {
		return Order{}, err
	}

	return order, nil
}

func (m *Model) OrderProducts(iOrdID uint) ([]OrderProduct, error) {

	query := `SELECT
				p.iProdID,
				p.vUrlName,
				p.vName,
				cat.vName AS vCategoryName,
				pa.vValue AS vAttributeValue,
				c.vName AS vColorValue,
				od.iQty,
				od.fRate,
				od.fPayable,
				COALESCE(od.cGiftWrap, '') AS cGiftWrap
			FROM orders o
				JOIN orders_dat od ON o.iOrdID = od.iOrdID
				JOIN product p ON od.iProdID = p.iProdID
				JOIN prodcat cat ON p.iPCatID = cat.iPCatID
				LEFT JOIN orders_dat_attrib_assoc odaa ON od.iOrdDatID = odaa.iOrdDatID
				LEFT JOIN product_attrib pa ON COALESCE(odaa.iProdAttribID, 0) = pa.iProdAttribID
				LEFT JOIN orders_dat_color_assoc odca ON od.iOrdDatID = odca.iOrdDatID
				LEFT JOIN product_color_assoc pca ON COALESCE(odca.iPCID, 0) = pca.iPCID
				LEFT JOIN color c ON pca.iColorID = c.iColorID
			WHERE o.iOrdID = ?`

	products := []OrderProduct{}
	if err := m.Select(&products, query, iOrdID); err != nil {
		return nil, err
	}

	return products, nil
}

func (m *Model) OrderShipments(iOrdID uint) ([]OrderShipment, error) {

	query := `SELECT 
				iOrdShipID,
				iOrdID,
				os.iCourierID,
				c.vName AS vCourierName,
				c.vLink,
				vShipCode,
				dtUpdatedAt
			FROM order_shipment os
				LEFT JOIN courier c ON os.iCourierID = c.iCourierID
			WHERE iOrdID = ?`

	shipments := []OrderShipment{}
	if err := m.Select(&shipments, query, iOrdID); err != nil {
		return nil, err
	}

	return shipments, nil
}

var insQuery = `INSERT INTO order_shipment 
				(iOrdID, iCourierID, vShipCode)
			VALUES 
				(?, ?, ?)`

var updtQuery = `UPDATE order_shipment SET 
				iOrdID = ?, 
				iCourierID = ?, 
				vShipCode = ?
			WHERE iOrdShipID = ?`

func (m *Model) SetOrderShipments(shipments []OrderShipment) ([]OrderShipment, error) {

	tx, err := m.Beginx()
	if err != nil {
		return nil, err
	}

	for _, shipment := range shipments {

		if shipment.IOrdShipID == 0 {
			result, err := tx.Exec(insQuery,
				shipment.IOrdID,
				shipment.ICourierID,
				shipment.VShipCode,
			)
			if err != nil {
				return nil, err
			}

			iOrdShipId, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}
			shipment.IOrdShipID = uint32(iOrdShipId)
		} else {
			result, err := tx.Exec(updtQuery,
				shipment.IOrdID,
				shipment.ICourierID,
				shipment.VShipCode,
				shipment.IOrdShipID,
			)
			if err != nil {
				return nil, err
			}

			_, err = result.RowsAffected()
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return shipments, nil
}

func (m *Model) DeleteShipment(iOrdShipID uint) (int, error) {

	delQuery := `DELETE FROM order_shipment WHERE iOrdShipID = ?`
	result, err := m.Exec(delQuery, iOrdShipID)
	if err != nil {
		return 0, err
	}

	deletedCount, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(deletedCount), nil
}
