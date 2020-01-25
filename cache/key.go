package cache

// Get 返回 key 对应的值
func Get(key string) (string, error) {
	return rc.Get(key).Result()
}

// Set 设置 key 对应的值
func Set(key, value string) (string, error) {
	return rc.Set(key, value, 0).Result()
}
