package webRequest

import (
	"strings"
)

type WebHeaderTarget struct {
	Name  string
	Value string
}

func (h *WebHeaderTarget) ToString() string {
	return h.Name + ":" + h.Value
}

type WebHeader struct {
	headers []*WebHeaderTarget
}

func (h *WebHeader) Add(key string, value string) {
	key = strings.ToLower(key)
	h.Del(key)
	h.headers = append(h.headers, &WebHeaderTarget{
		Name:  key,
		Value: value,
	})
}
func (h *WebHeader) AddHead(head *WebHeaderTarget) {
	h.Del(head.Name)
	h.headers = append(h.headers, head)
}

func (h *WebHeader) Set(key string, value string) {
	hd := h.Get(key)
	if hd != nil {
		hd.Value = value
	} else {
		h.headers = append(h.headers, &WebHeaderTarget{
			Name:  strings.ToLower(key),
			Value: value,
		})
	}
}
func (h *WebHeader) SetHead(head *WebHeaderTarget) {
	key := strings.ToLower(head.Name)
	for index, hd := range h.headers {
		if hd.Name == key {
			h.headers[index] = head
			return
		}
	}
	h.headers = append(h.headers, head)
}

func (h *WebHeader) Get(key string) *WebHeaderTarget {
	key = strings.ToLower(key)
	for _, hd := range h.headers {
		if hd.Name == key {
			return hd
		}
	}
	return nil
}
func (h *WebHeader) GetValue(key string) string {
	hd := h.Get(key)
	if hd == nil {
		return ""
	}
	return hd.Value
}

func (h *WebHeader) Has(key string) bool {
	return h.Get(key) != nil
}

func (h *WebHeader) Del(key string) bool {
	key = strings.ToLower(key)
	for index, hd := range h.headers {
		if hd.Name == key {
			h.headers = append(h.headers[:index], h.headers[index+1:]...)
			return true
		}
	}
	return false
}

func (h *WebHeader) ToString() string {
	result := []string{}
	for _, hd := range h.headers {
		result = append(result, hd.ToString())
	}
	return strings.Join(result, "\n")
}

func (h *WebHeader) ToMap() map[string]string {
	result := map[string]string{}
	for _, hd := range h.headers {
		result[hd.Name] = hd.Value
	}
	return result
}

// 获取排序
func (h *WebHeader) GetOrder() []string {
	result := make([]string, 0)
	for _, hd := range h.headers {
		result = append(result, hd.Name)
	}
	return result
}
