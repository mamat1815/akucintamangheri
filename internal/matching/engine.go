package matching

import (
	"campus-lost-and-found/internal/models"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type NotificationService interface {
	CreateNotification(userID uuid.UUID, title, body, refType string, refID uuid.UUID) error
}

type MatchingEngine struct {
	NotifService NotificationService
}

func NewMatchingEngine(notifService NotificationService) *MatchingEngine {
	return &MatchingEngine{NotifService: notifService}
}

// Haversine distance calculation
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371e3 // metres
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	d := R * c // in metres
	return d
}

func (e *MatchingEngine) RunMatching(foundItem *models.Item, lostAssets []models.Asset) {
	// Filter lost assets where:
	// lost_mode = true (already filtered by caller)
	// category matches

	for _, asset := range lostAssets {
		if asset.CategoryID != foundItem.CategoryID {
			continue
		}

		// Location Distance Score
		// Assuming Asset has a "Last Known Location" or we use the FoundEvent location if it was scanned?
		// Wait, Asset table doesn't have location.
		// However, the prompt says "Score asset based on: location distance (< 500m)".
		// This implies Asset MUST have a location.
		// But the schema for Asset is: id, owner_id, category_id, description, private_image_url, lost_mode, qr_code_url.
		// It does NOT have location.
		// Maybe it uses the "Found Events" location if the asset was scanned?
		// Or maybe the User inputs a location when marking as lost?
		// The prompt says "POST /assets/:id/lost-mode". Maybe we should accept location there?
		// But the schema is fixed.
		// Let's look at "Found Events". "found_events" has "location_id".
		// If an asset was scanned and reported found, it has a location.
		// But "Matching Engine" is for "Finder-First" flow: "POST /items/found ... run Smart Matching Engine".
		// This matches a NEWLY FOUND item against LOST assets.
		// If the asset is just "Lost" without location, we can't match by distance.
		// UNLESS the "Lost Mode" implies we know where it was lost?
		// Or maybe we match against "Found Events" of that asset?
		// No, "Found Events" are when someone finds it and scans QR.
		// "Items" are when someone finds something WITHOUT QR (Finder-First).
		// So we are matching Found Item (with location) against Lost Asset (without location?).
		// This is a gap.
		// However, I must follow the prompt.
		// "Score asset based on: location distance (< 500m)".
		// I will assume for now that we skip location score if asset has no location, OR
		// maybe I missed something.
		// Ah, "campus_locations" has lat/long.
		// Maybe the Asset has a "Home Location" or "Last Scanned Location"?
		// Let's assume we can't score location if we don't know it.
		// BUT, if the user reported it lost, maybe they provided a location?
		// The prompt for "PUT /assets/:id/lost-mode" doesn't specify body.
		// I'll assume for now we only match on Category and Time (if we had time of loss).
		// Wait, "time difference (< 24h)". Asset has `updated_at` which could be when it was marked lost.
		// So we can use `updated_at` of lost mode vs `created_at` of found item.
		
		// Let's check if I can add location to Asset or if I should just use what I have.
		// I will use `updated_at` for time score.
		// For location, since I cannot change schema, I will skip it or assume 0 distance if unknown (which is bad).
		// Actually, maybe the "Found Events" are relevant?
		// If an asset was scanned (FoundEvent), it has a location.
		// If I find an item, maybe it matches a FoundEvent?
		// No, FoundEvent is "I found this QR code". Item is "I found this object (no QR)".
		
		// Let's assume the prompt implies we should have location for lost assets.
		// But since I can't change schema, I will implement logic that checks if we have any location info.
		// Maybe I can check if there are any "FoundEvents" for this asset?
		// If there are, use the latest one's location?
		// That makes sense: Asset was scanned at Loc A (FoundEvent), then someone reports Found Item at Loc A.
		// But if it was scanned, the owner is already notified.
		
		// Let's just implement Time Score and Category Match.
		// And if I can't do location, I'll note it.
		// Wait, `campus_locations` table exists.
		// Maybe `Asset` should have `LocationID`? No.
		// I will stick to Category and Time.
		// And I'll add a dummy "Distance Score" that is always 100% if we can't calculate, or 0.
		// Actually, let's look at `Item` (Finder-First). It has `LocationID`.
		// If `Asset` has no location, we can't match distance.
		// I will implement the score logic but comment about the missing location on Asset.
		
		// Time Score
		// Difference between Found Item CreatedAt and Asset UpdatedAt (Lost Mode Time)
		timeDiff := foundItem.CreatedAt.Sub(asset.UpdatedAt).Hours()
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}
		
		timeScore := 0.0
		if timeDiff < 24 {
			timeScore = 100.0
		} else if timeDiff < 48 {
			timeScore = 50.0
		}

		// Final Score (Weighted)
		// If we had location, it would be 50% Loc + 50% Time.
		// Without location, let's just use Time?
		// Or maybe the prompt implies we SHOULD have location in Asset?
		// "assets (Owner-First) ... category_id ... description ... lost_mode ..."
		// No location.
		// I will proceed with Time Score only for now, but threshold is 80%.
		// So if Time < 24h, score is 100 >= 80 -> Match.
		
		finalScore := timeScore

		if finalScore >= 80 {
			// Notify Owner
			e.NotifService.CreateNotification(
				asset.OwnerID,
				"Potential Match Found!",
				fmt.Sprintf("An item matching your lost asset '%s' was reported found.", asset.Description),
				"POTENTIAL_MATCH",
				foundItem.ID,
			)
		}
	}
}
