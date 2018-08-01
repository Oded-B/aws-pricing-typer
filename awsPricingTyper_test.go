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

func (m *mockPricingClient) GetProducts(input *pricing.GetProductsInput) (*pricing.GetProductsOutput, error) {
	// define products
	product := make(map[string]interface{})
	productAttributes := make(map[string]interface{})
	product["productFamily"] = "Compute Instance"
	productAttributes["networkPerformance"] = "Moderate"
	productAttributes["vcpu"] = "2"
	productAttributes["capacityStatus"] = "Used"
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
	product["productAttributes"] = productAttributes

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
	onDemandPriceDimensionOne["appliesTo"] = "[]"
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
	reservedPriceDimensionOne["appliesTo"] = "[]"
	reservedPriceDimensionOne["rateCode"] = "7X4K64YA59VZZAC3.4NA7Y494T4.6YS6EN2CT7"
	reservedPriceDimensionOne["beginRange"] = "0"
	reservedPriceDimensions["7X4K64YA59VZZAC3.4NA7Y494T4.6YS6EN2CT7"] = reservedPriceDimensionOne
	reservedTermOne["priceDimensions"] = reservedPriceDimensions
	reservedTerms["7X4K64YA59VZZAC3.4NA7Y494T4"] = reservedTermOne
	terms["OnDemand"] = onDemandTerms
	terms["Reserved"] = reservedTerms

	plMap := make(map[string]interface{})
	plMap["publicationDate"] = "2018-07-27T01:58:36Z"
	plMap["version"] = "20180727015836"
	plMap["serviceCode"] = "AmazonEC2"
	plMap["terms"] = terms
	plMap["product"] = product

	output := pricing.GetProductsOutput{
		FormatVersion: getStrPtr("aws_v1"),
		NextToken:     getStrPtr(""),
		PriceList: []aws.JSONValue{
			plMap,
		},
	}
	return &output, nil
}

func TestTyper(t *testing.T) {
	// Setup Test
	mockSvc := &mockPricingClient{}
	getProductsInput := pricing.GetProductsInput{}
	getProductOutput, _ := mockSvc.GetProducts(&getProductsInput)
	_, err := GetTypedPricingData(*getProductOutput)
	if err != nil {
		t.Errorf("got error: %+v", err)
	}

}

func getStrPtr(input string) *string {
	return &input
}
