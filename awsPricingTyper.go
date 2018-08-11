package awsPricingTyper

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/service/pricing"
)

// GetTypedPricingData takes the raw output from the AWS API and returns typed data in structs
func GetTypedPricingData(getProductsOutput pricing.GetProductsOutput) (pricingData []PricingDocument, err error) {
	type JSONValue map[string]interface{}
	for _, item := range getProductsOutput.PriceList {
		var pDoc PricingDocument
		for k, v := range item {
			switch val := v.(type) {
			case string:
				switch k {
				case "publicationDate":
					pDoc.PublicationDate = val
				case "version":
					pDoc.Version = val
				case "serviceCode":
					pDoc.ServiceCode = val
				default:
					err = fmt.Errorf("unexpected price list item: %+v", k)
					return
				}
			case map[string]interface{}:
				switch k {
				case "product":
					var result product
					result, err = processProduct(v)
					if err != nil {
						return nil, err
					}
					pDoc.Products = append(pDoc.Products, result)
				case "terms":
					proTermsErr := processTerms(&pDoc, v)
					if proTermsErr != nil {
						err = fmt.Errorf("failed to process terms: %+v", proTermsErr)
						return nil, err
					}
				default:
					return nil, fmt.Errorf("unexpected price list item: %+v", k)
				}
			default:
				return nil, fmt.Errorf("unexpected type: %+v", val)
			}
		}
		pricingData = append(pricingData, pDoc)
	}
	return pricingData, err
}

func processProduct(v interface{}) (newProduct product, err error) {
	for k1, v1 := range v.(map[string]interface{}) {
		switch val := v1.(type) {
		case string:
			switch k1 {
			case "productFamily":
				if val != "Compute Instance" {
					continue
				}
				newProduct.ProductFamily = val
			case "sku":
				newProduct.SKU = val
			default:
				err = fmt.Errorf("unexpected field: %+v", k1)
			}
		case map[string]interface{}:
			for k2, v2 := range v1.(map[string]interface{}) {
				switch val := v2.(type) {
				case string:
					switch k2 {
					case "networkPerformance":
						newProduct.Attributes.NetworkPerformance = val
					case "vcpu":
						newProduct.Attributes.VCPU = val
					case "capacitystatus":
						newProduct.Attributes.CapacityStatus = val
					case "operatingSystem":
						newProduct.Attributes.OperatingSystem = val
					case "physicalProcessor":
						newProduct.Attributes.PhysicalProcessor = val
					case "ecu":
						newProduct.Attributes.ECU = val
					case "preInstalledSw":
						newProduct.Attributes.PreInstalledSw = val
					case "processorArchitecture":
						newProduct.Attributes.ProcessorArchitecture = val
					case "enhancedNetworkingSupported":
						newProduct.Attributes.EnhancedNetworkingSupported = val
					case "storage":
						newProduct.Attributes.Storage = val
					case "clockSpeed":
						newProduct.Attributes.ClockSpeed = val
					case "tenancy":
						newProduct.Attributes.Tenancy = val
					case "licenseModel":
						newProduct.Attributes.LicenseModel = val
					case "servicecode":
						newProduct.Attributes.ServiceCode = val
					case "currentGeneration":
						newProduct.Attributes.CurrentGeneration = val
					case "dedicatedEbsThroughput":
						newProduct.Attributes.DedicatedEbsThroughput = val
					case "servicename":
						newProduct.Attributes.ServiceName = val
					case "instanceType":
						newProduct.Attributes.InstanceType = val
					case "normalizationSizeFactor":
						newProduct.Attributes.NormalizationSizeFactor = val
					case "processorFeatures":
						newProduct.Attributes.ProcessorFeatures = val
					case "operation":
						newProduct.Attributes.Operation = val
					case "memory":
						newProduct.Attributes.Memory = val
					case "locationType":
						newProduct.Attributes.LocationType = val
					case "instanceFamily":
						newProduct.Attributes.InstanceFamily = val
					case "usagetype":
						newProduct.Attributes.UsageType = val
					case "location":
						newProduct.Attributes.Location = val
					default:
						err = fmt.Errorf("unexpected attribute: %+v", k2)
						return
					}
				}
			}
		default:
			err = fmt.Errorf("bad type: %+v", val)
			return
		}
	}
	return
}

func processReservedTerms(v1 interface{}) (reservedTerms map[string]reservedTerm, err error) {
	reservedTerms = make(map[string]reservedTerm)
	for k2, v2 := range v1.(map[string]interface{}) {
		var newReservedTerm reservedTerm
		switch v2.(type) {
		case string:
			err = fmt.Errorf("unexpected item: %+v %+v", k2, v2)
			return
		default:
			for k3, v3 := range v2.(map[string]interface{}) {
				switch val := v3.(type) {
				case string:
					switch k3 {
					case "sku":
						newReservedTerm.sku = val
					case "offerTermCode":
						newReservedTerm.offerTermCode = val
					case "effectiveDate":
						newReservedTerm.effectiveDate = val
					}
				default:
					if k3 == "termAttributes" {
						for k3ta, v3ta := range v3.(map[string]interface{}) {
							switch k3ta {
							case "LeaseContractLength":
								newReservedTerm.termAttributes.LeaseContractLength = v3ta.(string)
							case "OfferingClass":
								newReservedTerm.termAttributes.OfferingClass = v3ta.(string)
							case "PurchaseOption":
								newReservedTerm.termAttributes.PurchaseOption = v3ta.(string)
							}
						}
					} else if k3 == "priceDimensions" {
						var newPriceDimensions []priceDimension
						for pdK, pdV := range v3.(map[string]interface{}) {
							newPriceDimension := priceDimension{}
							switch val := pdV.(type) {
							default:
								err = fmt.Errorf("got unexpected price dimension value: %+v", val)
							case map[string]interface{}:
								var newPDItem priceDimensionItem
								for pdiK, pdiV := range pdV.(map[string]interface{}) {
									switch pdiK {
									default:
										err = fmt.Errorf("got unexpected price dimension field: %+v", pdiK)
									case "unit":
										newPDItem.unit = pdiV.(string)
									case "pricePerUnit":
										for pdiKu, pdiKv := range pdiV.(map[string]interface{}) {
											pricePerUnit := make(map[string]float64)
											pdiKvStr := pdiKv.(string)
											pdiKvFloat, conErr := strconv.ParseFloat(pdiKvStr, 64)
											if conErr != nil {
												return nil, conErr
											}
											pricePerUnit[pdiKu] = pdiKvFloat
											newPDItem.pricePerUnit = append(newPDItem.pricePerUnit, pricePerUnit)
										}
									case "appliesTo":
										switch pdiV.(type) {
										case map[string]interface{}:
											// TODO: work out what to do with it
										case []interface{}:
											// TODO: work out what to do with it
										default:
											err = fmt.Errorf("unexpected type for appliesTo: %+v", pdiV)
											return
										}
									case "endRange":
										newPDItem.endRange = pdiV.(string)
									case "description":
										newPDItem.description = pdiV.(string)
									case "rateCode":
										newPDItem.rateCode = pdiV.(string)
									case "beginRange":
										newPDItem.beginRange = pdiV.(string)
									}
									newPriceDimension[pdK] = newPDItem
								}
								newPriceDimensions = append(newPriceDimensions, newPriceDimension)
							}
							newReservedTerm.priceDimensions = newPriceDimensions

						}
					}
				}
			}
		}

		reservedTerms[k2] = newReservedTerm
	}
	return
}

func processOnDemandTerms(v1 interface{}) (onDemandTerms map[string]onDemandTerm, err error) {
	onDemandTerms = make(map[string]onDemandTerm)
	for k2, v2 := range v1.(map[string]interface{}) {
		var newOnDemandTerm onDemandTerm
		switch v2.(type) {
		case string:
			err = fmt.Errorf("unexpected item: %+v %+v", k2, v2)
		default:
			for k3, v3 := range v2.(map[string]interface{}) {
				switch val := v3.(type) {
				case string:
					switch k3 {
					case "sku":
						newOnDemandTerm.sku = val
					case "offerTermCode":
						newOnDemandTerm.offerTermCode = val
					case "effectiveDate":
						newOnDemandTerm.effectiveDate = val
					}
				default:
					if k3 == "termAttributes" {
						if len(v3.(map[string]interface{})) > 0 {
							err = fmt.Errorf("unexpected term attributes for OnDemand: %+v", val)
						}
					} else if k3 == "priceDimensions" {
						var newPriceDimensions []priceDimension
						for pdK, pdV := range v3.(map[string]interface{}) {
							newPriceDimension := priceDimension{}
							switch val := pdV.(type) {
							default:
								err = fmt.Errorf("got unexpected price dimension value: %+v", val)
							case map[string]interface{}:
								var newPDItem priceDimensionItem
								for pdiK, pdiV := range pdV.(map[string]interface{}) {
									switch pdiK {
									default:
										err = fmt.Errorf("got unexpected price dimension field: %+v", pdiK)
									case "unit":
										newPDItem.unit = pdiV.(string)
									case "pricePerUnit":
										for pdiKu, pdiKv := range pdiV.(map[string]interface{}) {
											pricePerUnit := make(map[string]float64)
											pdiKvStr := pdiKv.(string)
											pdiKvFloat, conErr := strconv.ParseFloat(pdiKvStr, 64)
											if conErr != nil {
												return nil, conErr
											}
											pricePerUnit[pdiKu] = pdiKvFloat
											newPDItem.pricePerUnit = append(newPDItem.pricePerUnit, pricePerUnit)
										}
									case "endRange":
										newPDItem.endRange = pdiV.(string)
									case "description":
										newPDItem.description = pdiV.(string)
									case "rateCode":
										newPDItem.rateCode = pdiV.(string)
									case "beginRange":
										newPDItem.beginRange = pdiV.(string)
									case "appliesTo":
										switch pdiV.(type) {
										case []interface{}:
											// TODO: work out what to do with it
										default:
											err = fmt.Errorf("unexpected type for appliesTo: %+v", pdiV)
											return nil, err
										}

									}
									newPriceDimension[pdK] = newPDItem
								}
								newPriceDimensions = append(newPriceDimensions, newPriceDimension)
							}
							newOnDemandTerm.priceDimensions = newPriceDimensions

						}
					}
				}
			}
		}
		onDemandTerms[k2] = newOnDemandTerm
	}
	return
}

func processTerms(doc *PricingDocument, v interface{}) error {
	for k1, v1 := range v.(map[string]interface{}) {
		switch v1.(type) {
		case map[string]interface{}:
			switch k1 {
			case "OnDemand":
				result, err := processOnDemandTerms(v1)
				if err != nil {
					return err
				}
				doc.Terms.OnDemand = result
			case "Reserved":
				result, err := processReservedTerms(v1)
				if err != nil {
					return err
				}
				doc.Terms.Reserved = result

			}
		default:
			return fmt.Errorf("did not expect value: %+v", v1)
		}
	}
	return nil
}

//type onDemandTermAttributes struct {
//	// empty
//}

type onDemandTerm struct {
	sku           string
	effectiveDate string
	offerTermCode string
	//termAttributes  onDemandTermAttributes
	priceDimensions []priceDimension
}

type pricePerUnit map[string]float64

type priceDimensionItem struct {
	pricePerUnit []pricePerUnit
	unit         string
	endRange     string
	description  string
	rateCode     string
	beginRange   string
	//appliesTo    interface{}
}

type priceDimension map[string]priceDimensionItem

type reservedTerm struct {
	sku            string
	effectiveDate  string
	offerTermCode  string
	termAttributes struct {
		LeaseContractLength string
		OfferingClass       string
		PurchaseOption      string
	}
	priceDimensions []priceDimension
}

type product struct {
	ProductFamily string
	SKU           string
	Attributes    struct {
		NetworkPerformance          string
		VCPU                        string
		CapacityStatus              string
		OperatingSystem             string
		PhysicalProcessor           string
		ECU                         string
		PreInstalledSw              string
		ProcessorArchitecture       string
		EnhancedNetworkingSupported string
		Storage                     string
		ClockSpeed                  string
		Tenancy                     string
		LicenseModel                string
		ServiceCode                 string
		CurrentGeneration           string
		DedicatedEbsThroughput      string
		ServiceName                 string
		InstanceType                string
		NormalizationSizeFactor     string
		ProcessorFeatures           string
		Operation                   string
		Memory                      string
		LocationType                string
		InstanceFamily              string
		UsageType                   string
		Location                    string
	}
}

// PricingDocument is a structure for each of the returned slice items
// representing each resulting product and it's accompanying pricing detail
type PricingDocument struct {
	PublicationDate string
	SKU             string
	ServiceCode     string
	Version         string
	Products        []product
	Terms           struct {
		OnDemand map[string]onDemandTerm
		Reserved map[string]reservedTerm
	}
}
