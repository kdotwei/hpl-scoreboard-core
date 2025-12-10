package handler

import (
	"encoding/json"
	"net/http"
)

// ... (CreateScoreRequest Struct 省略)

func (h *Handler) CreateScore(w http.ResponseWriter, r *http.Request) {
	// ... (前面的邏輯省略) ...

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 4. 回傳成功 (201 Created) 與 JSON 結果
	w.WriteHeader(http.StatusCreated)

	// 👇 修正這裡：加上錯誤檢查
	if err := json.NewEncoder(w).Encode(score); err != nil {
		// 雖然 Header 已經寫出去了，但紀錄錯誤還是必要的
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
