# Far Horizons Go Conversion Checklist

## 1. Project Setup
- [X] Create Go module structure
- [X] Set up testing framework
- [X] Create initial package organization

## 2. Core Data Structures
- [X] Convert star_data_t to Go
- [X] Convert game state structures
- [X] Add strong typing for numeric fields (int16, int32, etc)
- [ ] Define interfaces for key components

## 3. Utility Functions
- [ ] Port distanceBetween()
- [ ] Port scan() function
- [ ] Create helper functions for common operations

## 4. Sqlite3 Implementation
- [ ] Design schema for game state
- [ ] Implement queries for data marsalling and unmarshalling
- [ ] Create serialization functions for game objects
- [ ] Implement deserialization functions
- [ ] Add validation for SQL input/output
- [ ] Add reports with pretty printing for human-readable output

## 5. Game Engine Core
- [ ] Port main game loop
- [ ] Convert game state management
- [ ] Implement turn processing
- [ ] Add state validation

## 6. Testing & Validation
- [ ] Create unit tests
- [ ] Add integration tests
- [ ] Verify game logic equivalence
- [ ] Test data format compatibility
- [ ] Add JSON format tests

## 7. Documentation
- [ ] Add godoc comments
- [ ] Create usage examples
- [ ] Document conversion decisions
- [ ] Document JSON schema

## 8. Performance Optimization
- [ ] Profile Go implementation
- [ ] Optimize critical paths
- [ ] Benchmark against C version
- [ ] Optimize JSON handling

## 9. Final Steps
- [ ] Clean up unused code
- [ ] Ensure idiomatic Go style
- [ ] Create build/run instructions
- [ ] Verify all JSON interactions
