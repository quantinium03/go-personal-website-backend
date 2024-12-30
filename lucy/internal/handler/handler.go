package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mini-projects/keylogger-server/internal/database"
	"github.com/mini-projects/keylogger-server/internal/tp"
	"github.com/mini-projects/keylogger-server/utils"
	"golang.org/x/crypto/bcrypt"
)

type CounterHandler struct {
	ApiCfg tp.ApiConf
}

type parameters struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cfg *CounterHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to generate hash")
		return
	}

	log.Println("Username:", params.Username)
	log.Println("Password:", params.Password)

	user, err := cfg.ApiCfg.DB.InsertUser(r.Context(), database.InsertUserParams{
		Username:  params.Username,
		Userhash:  string(hash),
		Updatedat: time.Now().UTC(),
		Createdat: time.Now().UTC(),
	})
	if err != nil {
		log.Println(err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, utils.DatabaseUserToUser(user))
}

func (cfg *CounterHandler) IncrementCounter(w http.ResponseWriter, r *http.Request) {
	count, err := cfg.ApiCfg.DB.GetCounter(r.Context())
	if err != nil {
		log.Printf("Error fetching counter: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "Couldn't fetch the counter from the database")
		return
	}

	params := tp.CounterParameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Printf("Invalid request body: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if params.Counter < 0 {
		utils.ResponseWithError(w, http.StatusBadRequest, "Counter cannot be negative")
		return
	}

	newCount, err := cfg.ApiCfg.DB.UpdateCounter(r.Context(), params.Counter+count)
	if err != nil {
		log.Printf("Error updating counter: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "Couldn't update the counter")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]int{"counter": int(newCount)})
}

func (cfg *CounterHandler) GetCounter(w http.ResponseWriter, r *http.Request) {
	count, err := cfg.ApiCfg.DB.GetCounter(r.Context())
	if err != nil {
		log.Printf("Error fetching counter: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "Couldn't fetch the counter from the database")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]int{"counter": int(count)})
}

func (cfg *CounterHandler) UpdateMouseStats(w http.ResponseWriter, r *http.Request) {
	mouse, err := cfg.ApiCfg.DB.GetMouseStats(r.Context())
	if err != nil {
		log.Printf("Error fetching mouse stats from db: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "Couldn't fetch mouse stats from the database")
		return
	}

	params := tp.MouseStats{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Printf("Invalid Request body: %v", err)
		utils.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if params.LeftClick < 0 || params.RightClick < 0 || params.MouseDistance < 0.0 {
		log.Printf("Stats can't be negative")
		utils.ResponseWithError(w, http.StatusBadRequest, "Mouse stats can't be negative")
		return
	}

	stats, err := cfg.ApiCfg.DB.UpdateMouseStats(r.Context(), database.UpdateMouseStatsParams{
		Leftclick:     int64(params.LeftClick) + mouse.Leftclick,
		Rightclick:    int64(params.RightClick) + mouse.Rightclick,
		Mousedistance: int64(params.MouseDistance) + mouse.Mousedistance,
	})
	if err != nil {
		log.Printf("Failed to update mouse stats: %v", err)
		utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to update the mouse stats")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]int {
		"leftClick": int(stats.Leftclick),
		"rightClick": int(stats.Rightclick),
		"mouseDistance": int(stats.Mousedistance),
	})
}

func (cfg *CounterHandler) GetMouseStats(w http.ResponseWriter, r *http.Request) {
	stats, err := cfg.ApiCfg.DB.GetMouseStats(r.Context());
	if err != nil {
		log.Printf("Error getting mouse stats from database: %v", err);
		utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to get mouse stats from database")
		return;
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]int {
		"leftClick": int(stats.Leftclick),
		"rightClick": int(stats.Rightclick),
		"mouseDistance": int(stats.Mousedistance),
	})
}
