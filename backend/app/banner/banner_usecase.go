package banner

type BannerUsecase struct {
	bannerModel BannerModel
}

func (u BannerUsecase) GetBanners() ([]*Banner, error) {
	banners, err := u.bannerModel.GetBanners()
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (u BannerUsecase) CreateBanner(imageURL string, clickURL string) (*Banner, error) {
	banner, err := u.bannerModel.CreateBanner(imageURL, clickURL)
	if err != nil {
		return nil, err
	}

	return banner, nil
}

func (u BannerUsecase) DeleteBannerById(bannerId string) (int64, error) {
	count, err := u.bannerModel.DeleteBannerById(bannerId)
	if err != nil {
		return 0, err
	}

	return count, nil
}
