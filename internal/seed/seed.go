package seed

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
	"github.com/minguyentt/traverse/internal/db"
	"github.com/minguyentt/traverse/internal/storage"
	"github.com/minguyentt/traverse/models"
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

var jobTitles = []string{
	"SwiftMedix24", "NurseNomad_88", "VitalPulseX", "DocDart77", "HealNestRX",
	"ScrubScout", "MedRoamer", "RN4Hire", "PatchByte", "SyringeSurfer",
	"ChartMaster3", "PulseTrackz", "BedpanBandit", "ICUHustler", "CodeBlueCrew",
	"DocOnWheels", "GauzeGhost", "ShiftDrifter", "ScalpelShark",
	"VitalsCrate",
}

var cities = []string{
	"Austin, TX", "Seattle, WA", "Denver, CO", "Atlanta, GA",
	"Phoenix, AZ", "Orlando, FL", "Nashville, TN", "Portland, OR",
	"Minneapolis, MN", "Kansas City, MO", "Columbus, OH", "Sacramento, CA",
	"Raleigh, NC", "Pittsburgh, PA", "Salt Lake City, UT", "Albuquerque, NM",
	"Boise, ID", "San Antonio, TX", "Richmond, VA", "Milwaukee, WI",
}

var (
	agency           = "Aya HealthCare"
	profession       = "Registered Nurse"
	assignmentLength = []string{
		"12 weeks", "13 weeks", "14 weeks",
	}
)
var experience = "1 year"

func Seed(s *storage.Storage, db *db.PGDB) {
	rand := rand.New(rand.NewSource(time.Now().Unix()))

	ctx := context.Background()
	users := generateUsers(30, rand)

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

func generateUsers(num int, r *rand.Rand) []*models.User {
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

func generateDummyContracts(num int, users []*models.User, r *rand.Rand) []*models.Contract {
	contracts := make([]*models.Contract, num)

	for i := range num {
		usr := users[r.Intn(len(users))]

		contracts[i] = &models.Contract{
			UserID:   usr.ID,
			JobTitle: jobTitles[r.Intn(len(jobTitles))],
			City:     cities[r.Intn(len(cities))],
			Agency:   agency,
			JobDetails: &models.ContractJobDetails{
				Profession:       profession,
				AssignmentLength: assignmentLength[r.Intn(len(assignmentLength))],
				Experience:       experience,
			},
		}
	}

	return contracts
}

// TODO: dummy reviews for seeding

// func generateDummyReviews(
// 	num int,
// 	users []*models.User,
// 	contracts []*models.Contract,
// 	r *rand.Rand,
// ) []*models.Review {
// 	return nil
// }
