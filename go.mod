module github.com/rasmus-rasmus/pokedex

go 1.21.4

replace web v0.0.0 => ./web

replace cache v0.0.0 => ./cache

require (
	cache v0.0.0
	web v0.0.0
)
