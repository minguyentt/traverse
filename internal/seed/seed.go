package seed

import (
	"context"
	"fmt"
	"log"
	"traverse/internal/db"
	"traverse/internal/storage"
	"traverse/models"
)

var usernames = []string{
	"SwiftMedix24", "NurseNomad_88", "VitalPulseX", "DocDart77",
	"HealNestRX", "ScrubScout", "MedRoamer", "RN4Hire", "PatchByte",
	"SyringeSurfer", "ChartMaster3", "PulseTrackz", "BedpanBandit", "ICUHustler",
	"CodeBlueCrew", "DocOnWheels", "GauzeGhost", "ShiftDrifter", "ScalpelShark",
	"VitalsCrate",
}

var passwords = []string{
	"cloudy789123", "sunny456321", "appletree123", "happycat7890",
	"mintyfresh1", "cooldog3344", "sleepyowl12", "greengrass9",
	"tinybear007", "bluesky1122", "funnyfrog88", "sweetcorn33", "lazyduck456",
	"yellowfox99", "sandybeach1", "luckyfish77", "redrose1234",
	"fastbunny22", "softwind890", "quietmoon55",
}

var contractNames = []string{
	"SwiftMedix24", "NurseNomad_88", "VitalPulseX", "DocDart77", "HealNestRX",
	"ScrubScout", "MedRoamer", "RN4Hire", "PatchByte", "SyringeSurfer",
	"ChartMaster3", "PulseTrackz", "BedpanBandit", "ICUHustler", "CodeBlueCrew",
	"DocOnWheels", "GauzeGhost", "ShiftDrifter", "ScalpelShark",
	"VitalsCrate",
}

var content = []string{
	"$38.50", "$41.25", "$36.75", "$44.00", "$39.80",
	"$42.15", "$35.60", "$46.20", "$40.00", "$43.75",
	"$37.90", "$45.00", "$39.00", "$41.95", "$48.10",
	"$36.25", "$47.35", "$43.00", "$38.85", "$49.50",
}

var cities = []string{
	"Austin, TX", "Seattle, WA", "Denver, CO", "Atlanta, GA",
	"Phoenix, AZ", "Orlando, FL", "Nashville, TN", "Portland, OR",
	"Minneapolis, MN", "Kansas City, MO", "Columbus, OH", "Sacramento, CA",
	"Raleigh, NC", "Pittsburgh, PA", "Salt Lake City, UT", "Albuquerque, NM",
	"Boise, ID", "San Antonio, TX", "Richmond, VA", "Milwaukee, WI",
}

var addresses = []string{
	"4212 Maple Ridge Ave", "903 Willow Glen Dr", "1780 Pine Hollow Rd", "3124 Crestview Blvd", "2206 Harbor Bend Ln",
	"6953 Westmont Parkway", "1412 Copperfield Ct", "3847 Oakshade Loop", "7503 Summit Heights Way",
	"2894 Ironwood Terrace", "5827 Meadowbrook Trail", "1178 Sandstone Dr",
	"4429 Evergreen Hollow Rd", "3096 Lantern Cove Ln",
	"2365 Ridgeway Point", "8945 Forest Knoll Rd", "7650 Silver Creek St",
	"1950 Brookstone Ct", "6243 Juniper Edge Pl", "3308 Golden Grove Blvd",
}

func Seed(s *storage.Storage, db *db.PGDB) {
	ctx := context.Background()
	users := generateUsers(30)

	tx, _ := db.Begin(ctx)
	for _, usr := range users {
		if err := s.Users.CreateUser(ctx, usr); err != nil {
			_ = tx.Rollback(ctx)
			log.Println("err creating user:", err)
			return
		}
	}

	tx.Commit(ctx)

	log.Println("Seeding completed")
}

func generateUsers(num int) []*models.User {
	users := make([]*models.User, num)

	// TODO: need to generate random passwords
	for i := range num {
		users[i] = &models.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}

		if err := users[i].Password.Set([]byte(passwords[i%len(passwords)] + fmt.Sprintf("%d", i))); err != nil {
			log.Println("error hashing the password:", err)
			return nil
		}
	}

	return users
}

// TODO: generate contracts
// TODO: generate reviews
