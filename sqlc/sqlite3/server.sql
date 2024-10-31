--  Copyright (c) 2024 Michael D Henderson. All rights reserved.


-- GetGameState returns the game state.
--
-- name: GetServerPaths :one
SELECT code, name
FROM game_state;