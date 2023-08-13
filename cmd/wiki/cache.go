package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	checksums = make(map[string]string)
)

// Checksum
func Checksum(src string) string {
	v, ok := checksums[src]
	if ok {
		return v
	}
	file, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		} else {
			log.Fatal(err)
		}
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// UpdateCache
func UpdateCache(cache *Cache, key string, typ FileType) {
	switch typ {
	case FileTypeObject:
		cache.Objects[key] = Checksum(fmt.Sprintf("./../../database/wiki/%s.yml", key))
	case FileTypeObjectLogo:
		cache.ObjectLogos[key] = Checksum(fmt.Sprintf("./../../database/wiki_logos/%s.jpg", key))
	case FileTypeDict:
		cache.Dicts[key] = Checksum(fmt.Sprintf("./../../database/dict/%s.yml", key))
	}
}

// LoadCache
func LoadCache() (cache *Cache, err error) {
	cacheBytes, err := os.ReadFile("./cache/cache.json")
	if err == nil {
		if err := json.Unmarshal(cacheBytes, &cache); err != nil {
			return nil, err
		}
		return cache, nil
	}
	if os.IsNotExist(err) {
		cache = &Cache{
			Objects:     make(map[string]string),
			ObjectLogos: make(map[string]string),
			Dicts:       make(map[string]string),
		}
		cacheBytes, err = json.Marshal(cache)
		if err != nil {
			return nil, err
		}
		if err := os.MkdirAll("./cache", 0644); err != nil {
			log.Fatalln(err)
		}
		if err := os.WriteFile("./cache/cache.json", cacheBytes, 0644); err != nil {
			return nil, err
		}
		return cache, nil
	}
	return nil, err
}

// SaveCache
func SaveCache(cache *Cache) error {
	b, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	return os.WriteFile("./cache/cache.json", b, 0644)
}

// IsChecksumChanged
func IsChecksumChanged(cache *Cache, key string, typ FileType) bool {
	switch typ {
	case FileTypeObject:
		v, ok := cache.Objects[key]
		if !ok {
			return true
		}
		return v != Checksum(fmt.Sprintf("./../../database/wiki/%s.yml", key))
	case FileTypeObjectLogo:
		v, ok := cache.ObjectLogos[key]
		if !ok {
			return true
		}
		return v != Checksum(fmt.Sprintf("./../../database/wiki_logos/%s.jpg", key))
	case FileTypeDict:
		v, ok := cache.Dicts[key]
		if !ok {
			return true
		}
		return v != Checksum(fmt.Sprintf("./../../database/dict/%s.yml", key))
	}
	return false
}
