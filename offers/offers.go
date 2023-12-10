package offers

import (
	"sort"
	"time"
)

type Category uint32

const (
	Restaurant Category = iota + 1
	Retail
	Hotel
	Activity
)

type Merchant struct {
	Id       uint32 `json:"id"`
	Name     string `json:"name"`
	Distance float64
}

type Offer struct {
	Id          uint32      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Category    Category    `json:"category"`
	Merchants   []*Merchant `json:"merchants"`
	ValidTo     string      `json:"valid_to"`
}

type Offers struct {
	Offers []*Offer `json:"offers"`
}

func getDateTimeFromString(dateString string) time.Time {
	layoutFormat := "2006-01-02"
	timeObject, err := time.Parse(layoutFormat, dateString)
	if err != nil {
		panic(err)
	}
	return timeObject
}

func (o *Offer) getNearestMerchant() {
	minDistanceMerchant := o.Merchants[0]
	numOfMerchants := len(o.Merchants)
	for i := 0; i < numOfMerchants; i++ {
		if o.Merchants[i].Distance < minDistanceMerchant.Distance {
			minDistanceMerchant = o.Merchants[i]
		}
	}
	// Replace list of merchants in offer with its nearest merchant
	o.Merchants = []*Merchant{minDistanceMerchant}
}

func (o *Offers) FilterOffers(checkInDate string, deltaDays int) [2]*Offer {
	//since we only return one offer in one category
	//we use array instead of hashmap to optimize storage and for faster sort
	chosenOffers := []*Offer{}

	//we use mapIndex to track whether a category has chosen offer or not
	mapIndex := map[Category]int{}
	checkInTime := getDateTimeFromString(checkInDate)
	validDate := checkInTime.Add(time.Duration(24*deltaDays) * time.Hour)
	for _, offer := range o.Offers {
		//check valid category
		if offer.Category == Hotel {
			continue
		}

		//check valid date
		lastDateValid := getDateTimeFromString(offer.ValidTo)
		if validDate.Unix() > lastDateValid.Unix() {
			continue
		}

		newOffer := offer
		newOffer.getNearestMerchant()
		// check if a category has offer or not
		if _, ok := mapIndex[newOffer.Category]; !ok {
			//if category doesn't have offer, append offer to the chosenOffer
			chosenOffers = append(chosenOffers, offer)
			index := len(mapIndex) - 1
			mapIndex[newOffer.Category] = index
		} else {
			index := mapIndex[newOffer.Category]
			existOffer := chosenOffers[index]
			//compare distance of newOffer and existOffer
			if existOffer.Merchants[0].Distance > newOffer.Merchants[0].Distance {
				chosenOffers[index] = newOffer
			}
		}
	}
	sort.Slice(chosenOffers, func(i, j int) bool {
		closerDistance := chosenOffers[i].Merchants[0].Distance
		greaterDistance := chosenOffers[j].Merchants[0].Distance
		return closerDistance < greaterDistance
	})

	return [2]*Offer{chosenOffers[0], chosenOffers[1]}
}
