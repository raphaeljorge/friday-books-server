package models

    import (
      "time"
    )

    type Book struct {
      ID                string    `json:"id"`
      Title             string    `json:"title"`
      Subtitle          string    `json:"subtitle"`
      Authors           []string  `json:"authors"`
      Editors           []string  `json:"editors"`
      Translators       []string  `json:"translators"`
      Publisher         string    `json:"publisher"`
      PublicationDate   time.Time `json:"publication_date"`
      ISBN10            string    `json:"isbn_10"`
      ISBN13            string    `json:"isbn_13"`
      Language          string    `json:"language"`
      Genres            []string  `json:"genres"`
      Tags              []string  `json:"tags"`
      Description       string    `json:"description"`
      Edition           string    `json:"edition"`
      PageCount         int       `json:"page_count"`
      TargetAudience    string    `json:"target_audience"`
      SeriesName        string    `json:"series_name"`
      VolumeNumber      int       `json:"volume_number"`
      Awards            []string  `json:"awards"`
      CountryOfOrigin   string    `json:"country_of_origin"`
      CoverImageURL     string    `json:"cover_image_url"`
      Price             float64   `json:"price"`
      Currency          string    `json:"currency"`
      Retailers         []string  `json:"retailers"`
      AverageRating     float64   `json:"average_rating"`
      RatingCount       int       `json:"rating_count"`
      Reviews           []Review  `json:"reviews"`
      RelatedISBNs      []string  `json:"related_isbns"`
      CustomNotes       string    `json:"custom_notes"`
      CreatedAt         time.Time `json:"created_at"`
      UpdatedAt         time.Time `json:"updated_at"`
      
      // Format-specific fields
      Format            string    `json:"format,omitempty"` // Physical book
      PrintDimensions   string    `json:"print_dimensions,omitempty"`
      Weight            float64   `json:"weight,omitempty"`
      Illustrators      []string  `json:"illustrators,omitempty"`
      IllustrationCount int       `json:"illustration_count,omitempty"`
      DustJacket        string    `json:"dust_jacket,omitempty"`
      PaperType         string    `json:"paper_type,omitempty"`
      
      // eBook fields
      EBookFormat       string    `json:"ebook_format,omitempty"`
      EBookSize         int64     `json:"ebook_size,omitempty"`
      
      // Audiobook fields
      Narrators         []string  `json:"narrators,omitempty"`
      Runtime           int       `json:"runtime,omitempty"`
      AudioFormat       string    `json:"audio_format,omitempty"`
      AudioSize         int64     `json:"audio_size,omitempty"`
      SampleAudioURL    string    `json:"sample_audio_url,omitempty"`
      ChapterIndexed    bool      `json:"chapter_indexed"`
      Bitrate           int       `json:"bitrate,omitempty"`
      DownloadURL       string    `json:"download_url,omitempty"`
      AudioPublisher    string    `json:"audio_publisher,omitempty"`
      AudioRightsOwner  string    `json:"audio_rights_owner,omitempty"`
      AudioReleaseDate  time.Time `json:"audio_release_date,omitempty"`
    }

    type Review struct {
      ID        string    `json:"id"`
      UserID    string    `json:"user_id"`
      Rating    int       `json:"rating"`
      Comment   string    `json:"comment"`
      CreatedAt time.Time `json:"created_at"`
    }
