package mw

import (
	"context"
	"net/http"
	"time"
)

func TimeLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts := make([]time.Time, 0)
		ts = append(ts, time.Now())
		ctx := context.WithValue(r.Context(), "timestamp", &ts)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		//ts = append(ts, time.Now())
		//for i := 1; i < len(ts); i++ {
		//	fmt.Print(ts[i].Sub(ts[i-1]).Microseconds(), " ")
		//}
		//fmt.Println(r.Method, " ", r.URL)

		//if duration > 15000 {
		//	fmt.Println(duration, " ", r.Method, " ", r.URL)
		//}
	})
}
