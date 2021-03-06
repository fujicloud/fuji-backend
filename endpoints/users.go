package endpoints

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "../models"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        if key, ok := r.Header["Authorization"]; ok {
            if key[0] != Api_key {
                w.WriteHeader(403)
                return
            }
        } else {
            w.WriteHeader(403)
            return
        }
        
        decoder := json.NewDecoder(r.Body)
        var user models.User
        decoder.Decode(&user)
        defer r.Body.Close()
        
        stmtIns, _ := Db.Prepare("INSERT INTO users (os, device, locale, voiceover, bold_text, reduce_motion, reduce_transparency) VALUES (?, ?, ?, ?, ?, ?, ?)")
        defer stmtIns.Close()
        
        insertResult, _ := stmtIns.Exec(user.Os, user.Device, user.Locale, user.Voiceover, user.Bold_text, user.Reduce_motion, user.Reduce_transparency)
        
        id, _ := insertResult.LastInsertId()
        result := Db.QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", id)
        
        var resultUser models.User
        result.Scan(&resultUser.Id, &resultUser.Os, &resultUser.Device, &resultUser.Locale, &resultUser.Voiceover, &resultUser.Bold_text, &resultUser.Reduce_motion, &resultUser.Reduce_transparency, &resultUser.Created_at)
        
        resultJson, _ := json.Marshal(resultUser)
        fmt.Fprintf(w, string(resultJson))
    default:
        w.WriteHeader(404)
    }
}
