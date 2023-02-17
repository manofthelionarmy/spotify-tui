package tui

import "spotify-tui/internal/models"

// SpotifySearchArtistsMsg is a message sent
type SpotifySearchArtistsMsg []*models.Artist

// SelectedArtistMsg is message sent when we've selected an artist
type SelectedArtistMsg *models.Artist

// ArtistsResponse is the artists response we sent as a message
type ArtistsResponse []*models.Artist

// SongsResponse is the songs response we sent as a message
type SongsResponse []*models.Song

// SpotifySearchSongsRespMsg is a message signaling we got back an artists songs from spotify api
type SpotifySearchSongsRespMsg []*models.Song

// SpotifyAlbumsResponse is the alumb response we sent as a message
type SpotifyAlbumsResponse []*models.Album
