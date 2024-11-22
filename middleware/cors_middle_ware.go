package middleware

import "net/http"

// CORSミドルウェア
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") // すべてのオリジンを許可
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // 許可するHTTPメソッド
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // 許可するヘッダー

        // OPTIONSメソッドに対しては204ステータスで即座に返す
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}
