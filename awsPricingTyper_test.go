package awsPricingTyper

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/aws/aws-sdk-go/service/pricing/pricingiface"
)

type mockPricingClient struct {
	pricingiface.PricingAPI
}

// flag which item to fail on for price list (top-level) of input
var mockPriceListFailure string

// flag which item to fail on for product of input
var mockProductFailure string

// flag which item to fail on for terms of input
var mockTermsFailure string

func getMockProduct() (output map[string]interface{}) {
	product := make(map[string]interface{})
	productAttributes := make(map[string]interface{})
	if mockProductFailure == "unimplementedProductFamily" {
		product["productFamily"] = "Bad Family"
	} else {
		product["productFamily"] = "Compute Instance"
	}
	if mockProductFailure == "unimplementedProductItem" {
		product["badItem"] = "Bad Value"
	}
	product["sku"] = "7X4K64YA59VZZAC3"
	if mockProductFailure == "unexpectedProductAttributeValueType" {
		badProductAttributeValue := make(map[string]interface{})
		productAttributes["networkPerformance"] = badProductAttributeValue
	}
	productAttributes["networkPerformance"] = "Moderate"
	productAttributes["vcpu"] = "2"
	productAttributes["capacitystatus"] = "Used"
	productAttributes["operatingSystem"] = "Linux"
	productAttributes["physicalProcessor"] = "Intel Xeon E5-2676 v3 (Haswell)"
	productAttributes["ecu"] = "6.5"
	productAttributes["preInstalledSw"] = "NA"
	productAttributes["processorArchitecture"] = "64-bit"
	productAttributes["enhancedNetworkingSupported"] = "Yes"
	productAttributes["storage"] = "EBS"
	productAttributes["clockSpeed"] = "2.4 GHz"
	productAttributes["tenancy"] = "Shared"
	productAttributes["licenseModel"] = "No License required"
	productAttributes["servicecode"] = "AmazonEC2"
	productAttributes["currentGeneration"] = "Yes"
	productAttributes["dedicatedEbsThroughput"] = "450 Mbps"
	productAttributes["servicename"] = "Amazon Elastic Compute Cloud"
	productAttributes["instanceType"] = "m4.large"
	productAttributes["normalizationSizeFactor"] = "4"
	productAttributes["processorFeatures"] = "Intel AVX; Intel AVX2; Intel Turbo"
	productAttributes["operation"] = "RunInstances"
	productAttributes["memory"] = "8 GiB"
	productAttributes["locationType"] = "AWS Region"
	productAttributes["instanceFamily"] = "General purpose"
	productAttributes["usagetype"] = "EU-BoxUsage:m4.large"
	productAttributes["location"] = "EU (Ireland)"
	if mockProductFailure == "unexpectedProductAttribute" {
		productAttributes["badAttr"] = "a value"
	}
	product["productAttributes"] = productAttributes
	return product
}

func getMockTerms() map[string]interface{} {

	var mockReservedAppliesTo interface{}
	if mockTermsFailure == "unexpectedTypeForAppliesToReserved" {
		mockReservedAppliesTo = "bad type"
	} else {
		mockReservedAppliesTo = make(map[string]interface{})
	}
	var mockOnDemandAppliesTo interface{}
	if mockTermsFailure == "unexpectedTypeForAppliesToOnDemand" {
		mockOnDemandAppliesTo = "bad type"
	} else {
		mockOnDemandAppliesTo = make(map[string]interface{})
	}

	terms := make(map[string]interface{})
	onDemandTerms := make(map[string]interface{})
	onDemandTermOneAttrs := make(map[string]interface{})
	onDemandTermOne := make(map[string]interface{})
	onDemandTermOne["sku"] = "7X4K64YA59VZZAC3"
	onDemandTermOne["effectiveDate"] = "2018-07-01T00:00:00Z"
	onDemandTermOne["offerTermCode"] = "JRTCKXETXF"
	onDemandTermOne["termAttributes"] = onDemandTermOneAttrs
	onDemandPriceDimensionOne := make(map[string]interface{})
	onDemandPriceDimensionOnePricePerUnit := make(map[string]interface{})
	onDemandPriceDimensionOnePricePerUnit["USD"] = "0.1110000000"
	onDemandPriceDimensionOne["pricePerUnit"] = onDemandPriceDimensionOnePricePerUnit
	onDemandPriceDimensions := make(map[string]interface{})
	onDemandPriceDimensionOne["unit"] = "Hrs"
	onDemandPriceDimensionOne["endRange"] = "Inf"
	onDemandPriceDimensionOne["description"] = "$0.111 per On Demand Linux m4.large Instance Hour"
	onDemandPriceDimensionOne["endRange"] = "Inf"
	onDemandPriceDimensionOne["appliesTo"] = mockOnDemandAppliesTo
	onDemandPriceDimensionOne["rateCode"] = "7X4K64YA59VZZAC3.JRTCKXETXF.6YS6EN2CT7"
	onDemandPriceDimensionOne["beginRange"] = "0"
	onDemandPriceDimensions["ABCDEFGHIJK.LMNOPQRST.UVWXYZ"] = onDemandPriceDimensionOne
	onDemandTermOne["priceDimensions"] = onDemandPriceDimensions
	onDemandTerms["7X4K64YA59VZZAC3.JRTCKXETXF"] = onDemandTermOne

	reservedTerms := make(map[string]interface{})
	reservedTermOneAttrs := make(map[string]interface{})
	reservedTermOne := make(map[string]interface{})
	reservedTermOne["sku"] = "7X4K64YA59VZZAC3"
	reservedTermOne["effectiveDate"] = "2017-04-30T23:59:59Z"
	reservedTermOne["offerTermCode"] = "4NA7Y494T4"
	reservedTermOneAttrs["LeaseContractLength"] = "1yr"
	reservedTermOneAttrs["OfferingClass"] = "standard"
	reservedTermOneAttrs["PurchaseOption"] = "No Upfront"
	reservedTermOne["termAttributes"] = reservedTermOneAttrs
	reservedPriceDimensionOne := make(map[string]interface{})
	reservedPriceDimensionOnePricePerUnit := make(map[string]interface{})
	reservedPriceDimensionOnePricePerUnit["USD"] = "0.0756"
	reservedPriceDimensionOne["pricePerUnit"] = reservedPriceDimensionOnePricePerUnit
	reservedPriceDimensions := make(map[string]interface{})
	reservedPriceDimensionOne["unit"] = "Hrs"
	reservedPriceDimensionOne["endRange"] = "Inf"
	reservedPriceDimensionOne["description"] = "Linux/UNIX (Amazon VPC), m4.large reserved instance applied"
	reservedPriceDimensionOne["endRange"] = "Inf"
	reservedPriceDimensionOne["appliesTo"] = mockReservedAppliesTo
	reservedPriceDimensionOne["rateCode"] = "7X4K64YA59VZZAC3.4NA7Y494T4.6YS6EN2CT7"
	reservedPriceDimensionOne["beginRange"] = "0"
	reservedPriceDimensions["7X4K64YA59VZZAC3.4NA7Y494T4.6YS6EN2CT7"] = reservedPriceDimensionOne
	reservedTermOne["priceDimensions"] = reservedPriceDimensions
	reservedTerms["7X4K64YA59VZZAC3.4NA7Y494T4"] = reservedTermOne
	terms["OnDemand"] = onDemandTerms
	if mockTermsFailure == "unexpectedTypeForTerms" {
		terms["Reserved"] = "badTypeValue"
	} else {
		terms["Reserved"] = reservedTerms
	}
	return terms
}

func getMockPriceList(product, terms map[string]interface{}) map[string]interface{} {
	priceList := make(map[string]interface{})
	priceList["publicationDate"] = "2018-07-27T01:58:36Z"
	priceList["version"] = "20180727015836"
	priceList["serviceCode"] = "AmazonEC2"
	priceList["terms"] = terms
	priceList["product"] = product
	if mockPriceListFailure == "unimplementedPriceListItem" {
		priceList["badItem"] = "invalid"
	}
	return priceList
}

// working client
func (m *mockPricingClient) GetProducts(input *pricing.GetProductsInput) (*pricing.GetProductsOutput, error) {
	// get mock products
	var product map[string]interface{}
	product = getMockProduct()
	// get mock terms
	var terms map[string]interface{}
	terms = getMockTerms()
	// get mock pricelist
	var plMap map[string]interface{}
	plMap = getMockPriceList(product, terms)

	output := pricing.GetProductsOutput{
		FormatVersion: getStrPtr("aws_v1"),
		NextToken:     getStrPtr(""),
		PriceList: []aws.JSONValue{
			plMap,
		},
	}
	return &output, nil
}

// client failing with bad pricing document item
func TestTyperWithGoodData(t *testing.T) {
	// Setup Test
	mockProductFailure = "good"
	mockSvc := &mockPricingClient{}
	getProductsInput := pricing.GetProductsInput{}
	var getProductsOutput *pricing.GetProductsOutput
	getProductsOutput, getProductsErr := mockSvc.GetProducts(&getProductsInput)
	if getProductsErr != nil {
		t.Errorf("got unexpected error: %+v", getProductsErr)
	}

	_, getDataErr := GetTypedPricingData(*getProductsOutput)
	if getDataErr != nil {
		t.Errorf("got error: %+v", getDataErr)
	}

}

// client failing with unimplemented product family item
func TestTyperWithUnimplementedProductFamily(t *testing.T) {
	mockProductFailure = "unimplementedProductFamily"
	mockSvc := &mockPricingClient{}
	getProductsInput := pricing.GetProductsInput{}
	var getProductsOutput *pricing.GetProductsOutput
	getProductsOutput, getProductsErr := mockSvc.GetProducts(&getProductsInput)
	if getProductsErr != nil {
		t.Errorf("got unexpected error: %+v", getProductsErr)
	}
	_, getDataErr := GetTypedPricingData(*getProductsOutput)
	if getDataErr != nil {
		t.Errorf("got error: %+v", getDataErr)
	}

}

// client failing with unimplemented price list failure item
func TestTyperWithUnimplementedPriceListItem(t *testing.T) {
	mockPriceListFailure = "unimplementedPriceListItem"
	mockSvc := &mockPricingClient{}
	getProductsInput := pricing.GetProductsInput{}
	var getProductsOutput *pricing.GetProductsOutput
	getProductsOutput, getProductsErr := mockSvc.GetProducts(&getProductsInput)
	if getProductsErr != nil {
		t.Errorf("got unexpected error: %+v", getProductsErr)
	}
	_, getDataErr := GetTypedPricingData(*getProductsOutput)
	if getDataErr == nil {
		t.Errorf("expected bad price list item error: %+v", getDataErr)
	}

}

// client failing with unimplemented price list failure item
func TestTyperWithUnimplementedProductItem(t *testing.T) {
	mockProductFailure = "unimplementedProductItem"
	mockSvc := &mockPricingClient{}
	getProductsInput := pricing.GetProductsInput{}
	var getProductsOutput *pricing.GetProductsOutput
	getProductsOutput, getProductsErr := mockSvc.GetProducts(&getProductsInput)
	if getProductsErr != nil {
		t.Errorf("got unexpected error: %+v", getProductsErr)
	}
	_, getDataErr := GetTypedPricingData(*getProductsOutput)
	if getDataErr == nil {
		t.Errorf("expected unimplemented product item error: %+v", getDataErr)
	}

}

// client failing with unimplemented product attribute
func TestTyperWithUnexpectedProductAttribute(t *testing.T) {
	mockProductFailure = "unexpectedProductAttribute"
	mockSvc := &mockPricingClient{}
	getProductsInput := pricing.GetProductsInput{}
	var getProductsOutput *pricing.GetProductsOutput
	getProductsOutput, getProductsErr := mockSvc.GetProducts(&getProductsInput)
	if getProductsErr != nil {
		t.Errorf("got unexpected error: %+v", getProductsErr)
	}
	_, getDataErr := GetTypedPricingData(*getProductsOutput)
	if getDataErr == nil {
		t.Errorf("expected unimplemented product attribute error: %+v", getDataErr)
	}

}

// client failing with unimplemented product attribute type
func TestTermsWithInvalidType(t *testing.T) {
	mockTermsFailure = "unexpectedTypeForTerms"
	mockSvc := &mockPricingClient{}
	getProductsInput := pricing.GetProductsInput{}
	var getProductsOutput *pricing.GetProductsOutput
	getProductsOutput, getProductsErr := mockSvc.GetProducts(&getProductsInput)
	if getProductsErr != nil {
		t.Errorf("got unexpected error: %+v", getProductsErr)
	}
	_, getDataErr := GetTypedPricingData(*getProductsOutput)
	if getDataErr == nil {
		t.Errorf("expected unimplemented product attribute type error: %+v", getDataErr)
	}

}

// client failing with unexpected appliesTo type
func TestAppliesToWithInvalidTypeReserved(t *testing.T) {
	mockTermsFailure = "unexpectedTypeForAppliesToReserved"
	mockSvc := &mockPricingClient{}
	getProductsInput := pricing.GetProductsInput{}
	var getProductsOutput *pricing.GetProductsOutput
	getProductsOutput, getProductsErr := mockSvc.GetProducts(&getProductsInput)
	if getProductsErr != nil {
		t.Errorf("got unexpected error: %+v", getProductsErr)
	}
	_, getDataErr := GetTypedPricingData(*getProductsOutput)
	if getDataErr == nil {
		t.Errorf("expected unimplemented type for appliesTo error: %+v", getDataErr)
	}

}

// client failing with unexpected appliesTo type
func TestAppliesToWithInvalidTypeOnDemand(t *testing.T) {
	mockTermsFailure = "unexpectedTypeForAppliesToOnDemand"
	mockSvc := &mockPricingClient{}
	getProductsInput := pricing.GetProductsInput{}
	var getProductsOutput *pricing.GetProductsOutput
	getProductsOutput, getProductsErr := mockSvc.GetProducts(&getProductsInput)
	if getProductsErr != nil {
		t.Errorf("got unexpected error: %+v", getProductsErr)
	}
	_, getDataErr := GetTypedPricingData(*getProductsOutput)
	if getDataErr == nil {
		t.Errorf("expected unimplemented type for appliesTo error: %+v", getDataErr)
	}

}

func getStrPtr(input string) *string {
	return &input
}
