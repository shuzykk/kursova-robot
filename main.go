package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

type pristrii struct {
	Id            int `json:"id" db:"id"`
	Signal1       int `json:"Signal1" db:"Signal1"`
	Signal2       int `json:"signal2" db:"signal2"`
	Signal3       int `json:"signal3" db:"signal3"`
	Signal4       int `json:"signal4" db:"signal4"`
	Pristriinomer int `json:"Pristriinomer"`
	Id1           int `json:"id1"`
	Id2           int `json:"id2"`
}

func postgres() (*pgx.Conn, error) {
	c, e2 := pgx.Connect(context.Background(), "postgres://postgres:1122@db:5432/nit?sslmode=disable")
	if e2 != nil {
		fmt.Fprintln(os.Stderr, "Помилка в сегменті №2", e2)
		return nil, e2
	}
	return c, nil
}

func GET(w http.ResponseWriter, r *http.Request) {
	x, e3 := postgres()
	if e3 != nil {
		fmt.Fprintln(w, "Помилка в сегменті №3", e3)
		return
	}
	defer x.Close(context.Background())

	var unit1T []pristrii
	var unit2T []pristrii
	var unit3T []pristrii

	y1, _ := x.Query(context.Background(), "SELECT id, Signal1, signal2, signal3, signal4 FROM unit1 ORDER BY id ASC")
	for y1.Next() {
		var p pristrii
		e4 := y1.Scan(&p.Id, &p.Signal1, &p.Signal2, &p.Signal3, &p.Signal4)
		if e4 != nil {
			fmt.Fprintln(w, "Помилка в сегменті №4", e4)
			return
		}
		unit1T = append(unit1T, p)
	}
	fmt.Fprintln(w, "Пристрій №1")
	for _, p := range unit1T {
		fmt.Fprintf(w, "Номер Дня: %v, Сигнал 1: %v, Сигнал 2: %v, Сигнал 3: %v, Сигнал 4: %v\n", p.Id, p.Signal1, p.Signal2, p.Signal3, p.Signal4)
	}

	y2, _ := x.Query(context.Background(), "SELECT id, Signal1, signal2, signal3, signal4 FROM unit2 ORDER BY id ASC")
	for y2.Next() {
		var p pristrii
		e5 := y2.Scan(&p.Id, &p.Signal1, &p.Signal2, &p.Signal3, &p.Signal4)
		if e5 != nil {
			fmt.Fprintln(w, "Помилка в сегменті №5", e5)
			return
		}
		unit2T = append(unit2T, p)
	}
	fmt.Fprintln(w, "\nПристрій №2")
	for _, p := range unit2T {
		fmt.Fprintf(w, "Номер Дня: %v, Сигнал 1: %v, Сигнал 2: %v, Сигнал 3: %v, Сигнал 4: %v\n", p.Id, p.Signal1, p.Signal2, p.Signal3, p.Signal4)
	}

	y3, _ := x.Query(context.Background(), "SELECT id, Signal1, signal2, signal3, signal4 FROM unit3 ORDER BY id ASC")
	for y3.Next() {
		var p pristrii
		e6 := y3.Scan(&p.Id, &p.Signal1, &p.Signal2, &p.Signal3, &p.Signal4)
		if e6 != nil {
			fmt.Fprintln(w, "Помилка в сегменті №6", e6)
			return
		}
		unit3T = append(unit3T, p)
	}
	fmt.Fprintln(w, "\nПристрій №3")
	for _, p := range unit3T {
		fmt.Fprintf(w, "Номер Дня: %v, Сигнал 1: %v, Сигнал 2: %v, Сигнал 3: %v, Сигнал 4: %v\n", p.Id, p.Signal1, p.Signal2, p.Signal3, p.Signal4)
	}
}

func PUT(w http.ResponseWriter, r *http.Request) {
	danni, e8 := io.ReadAll(r.Body)
	if e8 != nil {
		fmt.Fprintln(w, "Помилка в сегменті №8", e8)
		return
	} else if len(danni) == 0 {
		fmt.Fprintln(w, "Немає данних")
		return
	}
	var pristr pristrii
	e9 := json.Unmarshal(danni, &pristr)
	if e9 != nil {
		fmt.Fprintln(w, "Помилка в сегменті №9", e9)
		return
	}
	x, e10 := postgres()
	if e10 != nil {
		fmt.Fprintln(w, "Помилка в сегменті №10", e10)
		return
	}
	defer x.Close(context.Background())

	if pristr.Pristriinomer == 1 {
		var proverka bool
		y5, ell := x.Query(context.Background(), "SELECT EXISTS(SELECT 1 FROM unit1 WHERE id=$1)", pristr.Id)
		if ell != nil {
			fmt.Fprintln(w, "Помилка в сегменті №11", ell)
			return
		}
		if y5.Next() {
			el2 := y5.Scan(&proverka)
			y5.Close()
			if el2 != nil {
				fmt.Fprintln(w, "Помилка в сегменті №12", el2)
				return
			} else if !proverka {
				fmt.Fprintln(w, "Немає результатів запиту")
				return
			}
			_, e5 := x.Exec(context.Background(), "UPDATE unit1 SET Signal1 = $1, signal2 = $2, signal3 = $3, signal4 = $4 WHERE id = $5", pristr.Signal1, pristr.Signal2, pristr.Signal3, pristr.Signal4, pristr.Id)
			fmt.Fprintln(w, "Оновлення строки", e5)
		}
	}
	if pristr.Pristriinomer == 2 {
		var proverka bool
		y6, e13 := x.Query(context.Background(), "SELECT EXISTS(SELECT 1 FROM unit2 WHERE id=$1) AND EXISTS (SELECT 1 FROM unit2 WHERE id=$2)", pristr.Id1, pristr.Id2)
		if e13 != nil {
			fmt.Fprintln(w, "Помилка в сегменті №13", e13)
			return
		}
		if y6.Next() {
			e14 := y6.Scan(&proverka)
			y6.Close()
			if e14 != nil {
				fmt.Fprintln(w, "Помилка в сегменті №14", e14)
				return
			} else if !proverka {
				fmt.Fprintln(w, "Немає результатів запиту")
				return
			}
			_, e5 := x.Exec(context.Background(), "DELETE FROM unit2 WHERE id IN ($1, $2)", pristr.Id1, pristr.Id2)
			fmt.Fprintln(w, "Видалення строки", e5)
		}

	}

	if pristr.Pristriinomer == 3 {
		var exists bool
		y, err := x.Query(context.Background(), "SELECT EXISTS(SELECT 1 FROM unit3 WHERE id=$1)", pristr.Id)
		if err != nil {
			fmt.Fprintln(w, "Помилка в сегменті №15", err)
			return
		}
		if y.Next() {
			err2 := y.Scan(&exists)
			y.Close()
			if err2 != nil {
				fmt.Fprintln(w, "Помилка в сегменті №16", err2)
				return
			} else if !exists {
				fmt.Fprintln(w, "Немає результатів запиту")
				return
			}
			_, err3 := x.Exec(context.Background(), "UPDATE unit3 SET signal1 = $1, signal2 = $2, signal3 = $3, signal4 = $4  WHERE id = $5", pristr.Signal1, pristr.Signal2, pristr.Signal3, pristr.Signal4, pristr.Id)
			fmt.Fprintln(w, "Оновлення строки ", err3)
		}
	}
}

func POST(w http.ResponseWriter, r *http.Request) {
	danni, e15 := io.ReadAll(r.Body)
	if e15 != nil {
		fmt.Fprintln(w, "Помилка в сегменті №17", e15)
		return
	} else if len(danni) == 0 {
		fmt.Fprintln(w, "Немає данних")
		return
	}
	var pristr pristrii
	e16 := json.Unmarshal(danni, &pristr)
	if e16 != nil {
		fmt.Fprintln(w, "Помилка в сегменті №18", e16)
		return
	}
	x, e17 := postgres()
	if e17 != nil {
		fmt.Fprintln(w, "Помилка в сегменті №19", e17)
		return
	}
	defer x.Close(context.Background())

	if pristr.Pristriinomer == 1 {
		_, e5 := x.Exec(context.Background(), "INSERT INTO unit1 (Signal1, signal2, signal3, signal4) VALUES($1, $2, $3, $4)", pristr.Signal1, pristr.Signal2, pristr.Signal3, pristr.Signal4)
		fmt.Fprintln(w, "Вставка значень у пристрій №1", e5)
	}
	if pristr.Pristriinomer == 2 {
		_, e5 := x.Exec(context.Background(), "INSERT INTO unit2 (Signal1, signal2, signal3, signal4) VALUES($1, $2, $3, $4)", pristr.Signal1, pristr.Signal2, pristr.Signal3, pristr.Signal4)
		fmt.Fprintln(w, "Вставка значень у пристрій №2", e5)
	}
	if pristr.Pristriinomer == 3 {
		_, e5 := x.Exec(context.Background(), "INSERT INTO unit3 (Signal1, signal2, signal3, signal4) VALUES($1, $2, $3, $4)", pristr.Signal1, pristr.Signal2, pristr.Signal3, pristr.Signal4)
		fmt.Fprintln(w, "Вставка значень у пристрій №3", e5)
	}

}

func rivensygnala(w http.ResponseWriter, r *http.Request) {
	x, err := postgres()
	if err != nil {
		fmt.Fprintln(w, "Помилка в сегменті №20", err)
		return
	}
	defer x.Close(context.Background())

	var unit1T []pristrii
	y1, _ := x.Query(context.Background(), "SELECT id, Signal1, signal2, signal3, signal4 FROM unit1 ORDER BY id ASC")
	for y1.Next() {
		var p pristrii
		e19 := y1.Scan(&p.Id, &p.Signal1, &p.Signal2, &p.Signal3, &p.Signal4)
		if e19 != nil {
			fmt.Fprintln(w, "Помилка в сегменті №21", e19)
			return
		}
		unit1T = append(unit1T, p)
	}
	plSignal1Suma := 0
	plsignal2Suma := 0
	plsignal3Suma := 0
	plsignal4Suma := 0
	for _, p := range unit1T {
		plSignal1Suma += p.Signal1
		plsignal2Suma += p.Signal2
		plsignal3Suma += p.Signal3
		plsignal4Suma += p.Signal4
	}
	fmt.Fprintln(w, "Пристрій №1")
	if len(unit1T) > 0 {
		if plSignal1Suma/len(unit1T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №1 на пристрої №1")
		} else if plSignal1Suma/len(unit1T) > 4 && plSignal1Suma/len(unit1T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №1 на пристрої №1")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №1 на пристрої №1")
		}
		if plsignal2Suma/len(unit1T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №2 на пристрої №1")
		} else if plsignal2Suma/len(unit1T) > 4 && plsignal2Suma/len(unit1T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №2 на пристрої №1")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №2 на пристрої №1")
		}
		if plsignal3Suma/len(unit1T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №3 на пристрої №1")
		} else if plsignal3Suma/len(unit1T) > 4 && plsignal3Suma/len(unit1T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №3 на пристрої №1")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №3 на пристрої №1")
		}
		if plsignal4Suma/len(unit1T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №4 на пристрої №1")
		} else if plsignal4Suma/len(unit1T) > 4 && plsignal4Suma/len(unit1T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №4 на пристрої №1")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №4 на пристрої №1")
		}
	}

	var unit2T []pristrii
	y2, _ := x.Query(context.Background(), "SELECT id, Signal1, signal2, signal3, signal4 FROM unit2 ORDER BY id ASC")
	for y2.Next() {
		var p pristrii
		e20 := y2.Scan(&p.Id, &p.Signal1, &p.Signal2, &p.Signal3, &p.Signal4)
		if e20 != nil {
			fmt.Fprintln(w, "Помилка в сегменті №22", e20)
			return
		}
		unit2T = append(unit2T, p)
	}
	p2Signal1Suma := 0
	p2signal2Suma := 0
	p2signal3Suma := 0
	p2signal4Suma := 0
	for _, p := range unit2T {
		p2Signal1Suma += p.Signal1
		p2signal2Suma += p.Signal2
		p2signal3Suma += p.Signal3
		p2signal4Suma += p.Signal4
	}
	fmt.Fprintln(w, "Пристрій №2")
	if len(unit2T) > 0 {
		if p2Signal1Suma/len(unit2T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №1 на пристрої №2")
		} else if p2Signal1Suma/len(unit2T) > 4 && p2Signal1Suma/len(unit2T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №1 на пристрої №2")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №1 на пристрої №2")
		}
		if p2signal2Suma/len(unit2T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №2 на пристрої №2")
		} else if p2signal2Suma/len(unit2T) > 4 && p2signal2Suma/len(unit2T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №2 на пристрої №2")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №2 на пристрої №2")
		}
		if p2signal3Suma/len(unit2T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №3 на пристрої №2")
		} else if p2signal3Suma/len(unit2T) > 4 && p2signal3Suma/len(unit2T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №3 на пристрої №2")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №3 на пристрої №2")
		}
		if p2signal4Suma/len(unit2T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №4 на пристрої №2")
		} else if p2signal4Suma/len(unit2T) > 4 && p2signal4Suma/len(unit2T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №4 на пристрої №2")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №4 на пристрої №2")
		}
	}

	var unit3T []pristrii
	y3, _ := x.Query(context.Background(), "SELECT id, Signal1, signal2, signal3, signal4 FROM unit3 ORDER BY id ASC")
	for y3.Next() {
		var p pristrii
		e21 := y3.Scan(&p.Id, &p.Signal1, &p.Signal2, &p.Signal3, &p.Signal4)
		if e21 != nil {
			fmt.Fprintln(w, "Помилка в сегменті №23", e21)
			return
		}
		unit3T = append(unit3T, p)
	}
	p3Signal1Suma := 0
	p3signal2Suma := 0
	p3signal3Suma := 0
	p3signal4Suma := 0
	for _, p := range unit3T {
		p3Signal1Suma += p.Signal1
		p3signal2Suma += p.Signal2
		p3signal3Suma += p.Signal3
		p3signal4Suma += p.Signal4
	}
	fmt.Fprintln(w, "Пристрій №3")
	if len(unit3T) > 0 {
		if p3Signal1Suma/len(unit3T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №1 на пристрої №3")
		} else if p3Signal1Suma/len(unit3T) > 4 && p3Signal1Suma/len(unit3T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №1 на пристрої №3")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №1 на пристрої №3")
		}
		if p3signal2Suma/len(unit3T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №2 на пристрої №3")
		} else if p3signal2Suma/len(unit3T) > 4 && p3signal2Suma/len(unit3T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №2 на пристрої №3")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №2 на пристрої №3")
		}
		if p3signal3Suma/len(unit3T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №3 на пристрої №3")
		} else if p3signal3Suma/len(unit3T) > 4 && p3signal3Suma/len(unit3T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №3 на пристрої №3")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №3 на пристрої №3")
		}
		if p3signal4Suma/len(unit3T) <= 4 {
			fmt.Fprintln(w, "Зелений рівень сигналу №4 на пристрої №3")
		} else if p3signal4Suma/len(unit3T) > 4 && p3signal4Suma/len(unit3T) < 7 {
			fmt.Fprintln(w, "Увага! Жовтий рівень сигналу №4 на пристрої №3")
		} else {
			fmt.Fprintln(w, "Увага! Червоний рівень сигналу №4 на пристрої №3")
		}
	}

}

func kyrsova(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GET(w, r)
	} else if r.Method == http.MethodPost {
		POST(w, r)
	} else if r.Method == http.MethodPut {
		PUT(w, r)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", kyrsova)
	http.HandleFunc("/run_stats", rivensygnala)
	e := http.ListenAndServe(":8080", nil)
	if e != nil {
		fmt.Println("Помилка головної функції", e)
	}
}
