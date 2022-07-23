package main 

import ( "log" 
	"net/http" 
	"time"
) 

var AppDir = http.Dir("webapp")

func frontendHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path 	
	log.Println(path)
	serveFile(&AppDir, path, true, w, r)
	return
}

func serveFile(dir *http.Dir, fileName string, sendIndex bool, w http.ResponseWriter, r *http.Request) bool { 
	file, err := dir.Open(fileName) 
	if err != nil {
		if sendIndex { 
			sendIndexFile(dir, w, r) 
			return true
		}
		return false
	}
	defer file.Close() 
	fileStat, err := file.Stat() 
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError) 
		return false
	}
	
	if fileStat.IsDir() { 
		if sendIndex { 
			sendIndexFile(dir, w, r) 
			return true
		}
		return false
	}
	
	http.ServeContent(w, r, fileName, fileStat.ModTime(), file) 
	return true
}

func sendIndexFile(dir *http.Dir, w http.ResponseWriter, r *http.Request) { 
	indexFile, err := dir.Open("index.html") 
	if err != nil {
		w.WriteHeader(http.StatusNotFound) 
		return
	}
	indexStat, err := indexFile.Stat() 
	if err != nil { 
		w.WriteHeader(http.StatusInternalServerError) 
		return
	}
	http.ServeContent(w, r, "index.html", indexStat.ModTime(), indexFile)
}

func main() { 
	s := &http.Server{ 
		Addr: "0.0.0.0:8080", 
		Handler: http.HandlerFunc(frontendHandler), 
		ReadTimeout: 10 * time.Second, 
		WriteTimeout: 10 * time.Second, 
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
