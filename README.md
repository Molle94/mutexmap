# mutexmap
MutexMap provides a generics-based go standard map protected by a RWMutex. It exposes the same methods as sync.Map but with better type safety and performance in scenarios with up to four cores.
For benchmarks with sync.Map compare (https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c)

## Usage
Create a new MutexMap with `string` keys and `*http.Client` values  
`m := NewMutexMap[string, *http.Client]() `
