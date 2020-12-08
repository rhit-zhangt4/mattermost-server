// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/disintegration/imaging"
	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/mattermost/mattermost-server/v5/utils"
)

const (
	MaxEmojiFileSize       = 1 << 20 // 1 MB
	MaxEmojiWidth          = 128
	MaxEmojiHeight         = 128
	MaxEmojiOriginalWidth  = 1028
	MaxEmojiOriginalHeight = 1028
)

func (a *App) CreateEmoji(sessionUserId string, emoji *model.Emoji, multiPartImageData *multipart.Form) (*model.Emoji, *model.AppError) {
	if !*a.Config().ServiceSettings.EnableCustomEmoji {
		return nil, model.NewAppError("UploadEmojiImage", "api.emoji.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	if len(*a.Config().FileSettings.DriverName) == 0 {
		return nil, model.NewAppError("GetEmoji", "api.emoji.storage.app_error", nil, "", http.StatusNotImplemented)
	}

	// wipe the emoji id so that existing emojis can't get overwritten
	emoji.Id = ""

	// do our best to validate the emoji before committing anything to the DB so that we don't have to clean up
	// orphaned files left over when validation fails later on
	emoji.PreSave()
	if err := emoji.IsValid(); err != nil {
		return nil, err
	}

	if emoji.CreatorId != sessionUserId {
		return nil, model.NewAppError("createEmoji", "api.emoji.create.other_user.app_error", nil, "", http.StatusForbidden)
	}

	if existingEmoji, err := a.Srv().Store.Emoji().GetByName(emoji.Name, true); err == nil && existingEmoji != nil {
		return nil, model.NewAppError("createEmoji", "api.emoji.create.duplicate.app_error", nil, "", http.StatusBadRequest)
	}

	imageData := multiPartImageData.File["image"]
	if len(imageData) == 0 {
		err := model.NewAppError("Context", "api.context.invalid_body_param.app_error", map[string]interface{}{"Name": "createEmoji"}, "", http.StatusBadRequest)
		return nil, err
	}

	if err := a.UploadEmojiImage(emoji.Id, imageData[0]); err != nil {
		return nil, err
	}

	emoji, err := a.Srv().Store.Emoji().Save(emoji)
	if err != nil {
		return nil, model.NewAppError("CreateEmoji", "app.emoji.create.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}

	public_emoji := &model.PublicEmoji{
		EmojiId: emoji.Id,
	}

	_, err = a.Srv().Store.PublicEmoji().Save(public_emoji)
	if err != nil {
		return nil, model.NewAppError("CreatePublicEmoji", "app.emoji.create.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}

	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_EMOJI_ADDED, "", "", "", nil)
	message.Add("emoji", emoji.ToJson())
	a.Publish(message)
	return emoji, nil
}

func (a *App) CreatePrivateEmoji(sessionUserId string, emoji *model.Emoji, multiPartImageData *multipart.Form) (*model.Emoji, *model.AppError) {
	if !*a.Config().ServiceSettings.EnableCustomEmoji {
		return nil, model.NewAppError("UploadEmojiImage", "api.emoji.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	if len(*a.Config().FileSettings.DriverName) == 0 {
		return nil, model.NewAppError("GetEmoji", "api.emoji.storage.app_error", nil, "", http.StatusNotImplemented)
	}

	// wipe the emoji id so that existing emojis can't get overwritten
	emoji.Id = ""

	// do our best to validate the emoji before committing anything to the DB so that we don't have to clean up
	// orphaned files left over when validation fails later on
	emoji.PreSave()
	if err := emoji.IsValid(); err != nil {
		return nil, err
	}

	if emoji.CreatorId != sessionUserId {
		return nil, model.NewAppError("createEmoji", "api.emoji.create.other_user.app_error", nil, "", http.StatusForbidden)
	}

	if existingEmoji, err := a.Srv().Store.Emoji().GetByName(emoji.Name, true); err == nil && existingEmoji != nil {
		return nil, model.NewAppError("createEmoji", "api.emoji.create.duplicate.app_error", nil, "", http.StatusBadRequest)
	}

	imageData := multiPartImageData.File["image"]
	if len(imageData) == 0 {
		err := model.NewAppError("Context", "api.context.invalid_body_param.app_error", map[string]interface{}{"Name": "createEmoji"}, "", http.StatusBadRequest)
		return nil, err
	}

	if err := a.UploadEmojiImage(emoji.Id, imageData[0]); err != nil {
		return nil, err
	}
	// TODO: set isPublic to false
	// TODO: change query
	emoji, err := a.Srv().Store.Emoji().Save(emoji)
	if err != nil {
		return nil, model.NewAppError("CreateEmoji", "app.emoji.create.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	emoji_access := &model.EmojiAccess{
		EmojiId: emoji.Id,
		UserId:  sessionUserId,
	}
	_, err = a.Srv().Store.EmojiAccess().Save(emoji_access)
	if err != nil {
		return nil, model.NewAppError("CreateEmoji", "app.emojiAceess.create.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_EMOJI_ADDED, "", "", "", nil)
	message.Add("emoji", emoji.ToJson())
	a.Publish(message)
	return emoji, nil
}

func (a *App) GetEmojiList(page, perPage int, sort string) ([]*model.Emoji, *model.AppError) {
	list, err := a.Srv().Store.Emoji().GetList(page*perPage, perPage, sort)
	if err != nil {
		return nil, model.NewAppError("GetEmojiList", "app.emoji.get_list.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return list, nil
}

func (a *App) GetPrivateEmojiList(page, perPage int, sort string, userid string) ([]*model.Emoji, *model.AppError) {
	// TODO: change query with userid
	//list, err := a.Srv().Store.Emoji().GetList(page*perPage, perPage, sort)
	list, err := a.Srv().Store.EmojiAccess().GetMultipleByUserId([]string{userid})
	if err != nil {
		return nil, model.NewAppError("GetEmojiList", "app.emoji.get_list.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	var emojis []*model.Emoji
	for _, v := range list {
		emoji, err := a.Srv().Store.Emoji().Get(v.EmojiId, true)
		if err != nil {
			return nil, model.NewAppError("GetEmojiPrivateList", "app.emoji.get_list.internal_error", nil, err.Error(), http.StatusInternalServerError)
		}

		emojis = append(emojis, emoji)

	}

	return emojis, nil
}

func (a *App) GetPublicEmojiList(page, perPage int, sort string) ([]*model.Emoji, *model.AppError) {
	// TODO: change query with userid
	//list, err := a.Srv().Store.Emoji().GetList(page*perPage, perPage, sort)
	list, err := a.Srv().Store.PublicEmoji().GetAllPublicEmojis()
	if err != nil {
		return nil, model.NewAppError("GetPublicEmoji", "app.emoji.get_list.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	var emojis []*model.Emoji
	for _, v := range list {
		emoji, err := a.Srv().Store.Emoji().Get(v.EmojiId, true)
		if err != nil {
			return nil, model.NewAppError("GetEmojiPublicList", "app.emoji.get_list.internal_error", nil, err.Error(), http.StatusInternalServerError)
		}

		emojis = append(emojis, emoji)
	}

	return emojis, nil
}

func (a *App) UploadEmojiImage(id string, imageData *multipart.FileHeader) *model.AppError {
	if !*a.Config().ServiceSettings.EnableCustomEmoji {
		return model.NewAppError("UploadEmojiImage", "api.emoji.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	if len(*a.Config().FileSettings.DriverName) == 0 {
		return model.NewAppError("UploadEmojiImage", "api.emoji.storage.app_error", nil, "", http.StatusNotImplemented)
	}

	file, err := imageData.Open()
	if err != nil {
		return model.NewAppError("uploadEmojiImage", "api.emoji.upload.open.app_error", nil, "", http.StatusBadRequest)
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)

	// make sure the file is an image and is within the required dimensions
	config, _, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return model.NewAppError("uploadEmojiImage", "api.emoji.upload.image.app_error", nil, "", http.StatusBadRequest)
	}

	if config.Width > MaxEmojiOriginalWidth || config.Height > MaxEmojiOriginalHeight {
		return model.NewAppError("uploadEmojiImage", "api.emoji.upload.large_image.too_large.app_error", map[string]interface{}{
			"MaxWidth":  MaxEmojiOriginalWidth,
			"MaxHeight": MaxEmojiOriginalHeight,
		}, "", http.StatusBadRequest)
	}

	if config.Width > MaxEmojiWidth || config.Height > MaxEmojiHeight {
		data := buf.Bytes()
		newbuf := bytes.NewBuffer(nil)
		info, err := model.GetInfoForBytes(imageData.Filename, data)
		if err != nil {
			return err
		}

		if info.MimeType == "image/gif" {
			gif_data, err := gif.DecodeAll(bytes.NewReader(data))
			if err != nil {
				return model.NewAppError("uploadEmojiImage", "api.emoji.upload.large_image.gif_decode_error", nil, "", http.StatusBadRequest)
			}

			resized_gif := resizeEmojiGif(gif_data)
			if err := gif.EncodeAll(newbuf, resized_gif); err != nil {
				return model.NewAppError("uploadEmojiImage", "api.emoji.upload.large_image.gif_encode_error", nil, "", http.StatusBadRequest)
			}

			if _, err := a.WriteFile(newbuf, getEmojiImagePath(id)); err != nil {
				return err
			}
		} else {
			img, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				return model.NewAppError("uploadEmojiImage", "api.emoji.upload.large_image.decode_error", nil, "", http.StatusBadRequest)
			}

			resized_image := resizeEmoji(img, config.Width, config.Height)
			if err := png.Encode(newbuf, resized_image); err != nil {
				return model.NewAppError("uploadEmojiImage", "api.emoji.upload.large_image.encode_error", nil, "", http.StatusBadRequest)
			}
			if _, err := a.WriteFile(newbuf, getEmojiImagePath(id)); err != nil {
				return err
			}
		}
	}

	_, appErr := a.WriteFile(buf, getEmojiImagePath(id))
	return appErr
}

func (a *App) DeletePrivateEmojiAccess(userid string, emojiId string) *model.AppError {
	if err := a.Srv().Store.EmojiAccess().DeleteAccessByUserIdAndEmojiId(userid, emojiId); err != nil {
		return model.NewAppError("DeleteEmoji", "app.emoji.delete.app_error", nil, "id="+emojiId+", err="+err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (a *App) DeleteEmojiWithAccess(userid string, emoji *model.Emoji) *model.AppError {
	if err := a.Srv().Store.PublicEmoji().DeleteAccessByEmojiId(emoji.Id); err != nil {
		return model.NewAppError("DeleteEmoji", "app.emoji.delete.app_error", nil, "id="+emoji.Id+", err="+err.Error(), http.StatusInternalServerError)
	}
	if err := a.Srv().Store.EmojiAccess().DeleteAccessByEmojiId(emoji.Id); err != nil {
		return model.NewAppError("DeleteEmoji", "app.emoji.delete.app_error", nil, "id="+emoji.Id+", err="+err.Error(), http.StatusInternalServerError)
	}
	if err := a.Srv().Store.Emoji().Delete(emoji, model.GetMillis()); err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return model.NewAppError("DeleteEmoji", "app.emoji.delete.no_results", nil, "id="+emoji.Id+", err="+err.Error(), http.StatusNotFound)
		default:
			return model.NewAppError("DeleteEmoji", "app.emoji.delete.app_error", nil, "id="+emoji.Id+", err="+err.Error(), http.StatusInternalServerError)
		}
	}

	a.deleteEmojiImage(emoji.Id)
	a.deleteReactionsForEmoji(emoji.Name)
	return nil
}

func (a *App) DeleteEmoji(emoji *model.Emoji) *model.AppError {
	if err := a.Srv().Store.Emoji().Delete(emoji, model.GetMillis()); err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return model.NewAppError("DeleteEmoji", "app.emoji.delete.no_results", nil, "id="+emoji.Id+", err="+err.Error(), http.StatusNotFound)
		default:
			return model.NewAppError("DeleteEmoji", "app.emoji.delete.app_error", nil, "id="+emoji.Id+", err="+err.Error(), http.StatusInternalServerError)
		}
	}

	a.deleteEmojiImage(emoji.Id)
	a.deleteReactionsForEmoji(emoji.Name)
	return nil
}

func (a *App) GetEmoji(emojiId string) (*model.Emoji, *model.AppError) {
	if !*a.Config().ServiceSettings.EnableCustomEmoji {
		return nil, model.NewAppError("GetEmoji", "api.emoji.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	if len(*a.Config().FileSettings.DriverName) == 0 {
		return nil, model.NewAppError("GetEmoji", "api.emoji.storage.app_error", nil, "", http.StatusNotImplemented)
	}

	emoji, err := a.Srv().Store.Emoji().Get(emojiId, false)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return emoji, model.NewAppError("GetEmoji", "app.emoji.get.no_result", nil, err.Error(), http.StatusNotFound)
		default:
			return emoji, model.NewAppError("GetEmoji", "app.emoji.get.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return emoji, nil
}

func (a *App) GetEmojiByName(emojiName string) (*model.Emoji, *model.AppError) {
	if !*a.Config().ServiceSettings.EnableCustomEmoji {
		return nil, model.NewAppError("GetEmojiByName", "api.emoji.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	if len(*a.Config().FileSettings.DriverName) == 0 {
		return nil, model.NewAppError("GetEmojiByName", "api.emoji.storage.app_error", nil, "", http.StatusNotImplemented)
	}

	emoji, err := a.Srv().Store.Emoji().GetByName(emojiName, true)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return emoji, model.NewAppError("GetEmojiByName", "app.emoji.get_by_name.no_result", nil, err.Error(), http.StatusNotFound)
		default:
			return emoji, model.NewAppError("GetEmojiByName", "app.emoji.get_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}

	return emoji, nil
}

func (a *App) GetMultipleEmojiByName(names []string) ([]*model.Emoji, *model.AppError) {
	if !*a.Config().ServiceSettings.EnableCustomEmoji {
		return nil, model.NewAppError("GetMultipleEmojiByName", "api.emoji.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	emoji, err := a.Srv().Store.Emoji().GetMultipleByName(names)
	if err != nil {
		return nil, model.NewAppError("GetMultipleEmojiByName", "app.emoji.get_by_name.app_error", nil, fmt.Sprintf("names=%v, %v", names, err.Error()), http.StatusInternalServerError)
	}

	return emoji, nil
}

func (a *App) GetEmojiImage(emojiId string) ([]byte, string, *model.AppError) {
	_, storeErr := a.Srv().Store.Emoji().Get(emojiId, true)
	if storeErr != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(storeErr, &nfErr):
			return nil, "", model.NewAppError("GetEmojiImage", "app.emoji.get.no_result", nil, storeErr.Error(), http.StatusNotFound)
		default:
			return nil, "", model.NewAppError("GetEmojiImage", "app.emoji.get.app_error", nil, storeErr.Error(), http.StatusInternalServerError)
		}
	}

	img, appErr := a.ReadFile(getEmojiImagePath(emojiId))
	if appErr != nil {
		return nil, "", model.NewAppError("getEmojiImage", "api.emoji.get_image.read.app_error", nil, appErr.Error(), http.StatusNotFound)
	}

	_, imageType, err := image.DecodeConfig(bytes.NewReader(img))
	if err != nil {
		return nil, "", model.NewAppError("getEmojiImage", "api.emoji.get_image.decode.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return img, imageType, nil
}
func (a *App) GetCanAccessPrivateEmojiImage(emojiId string, userId string) *model.AppError {
	_, accessErr := a.Srv().Store.EmojiAccess().GetByUserIdAndEmojiId(userId, emojiId)
	if accessErr != nil {
		return model.NewAppError("getEmojiImage", "api.emoji.get_image.read.app_error", nil, accessErr.Error(), http.StatusNotFound)
	}
	return nil
}

func (a *App) SavePrivateEmoji(emojiId string, userId string) *model.AppError {
	_, accessErr := a.Srv().Store.EmojiAccess().GetByUserIdAndEmojiId(userId, emojiId)
	if accessErr == nil {
		return model.NewAppError("createEmoji", "api.emoji.create.duplicate.app_error", nil, "", http.StatusBadRequest)
	}
	emoji_access := &model.EmojiAccess{
		EmojiId: emojiId,
		UserId:  userId,
	}
	_, err := a.Srv().Store.EmojiAccess().Save(emoji_access)
	if err != nil {
		// return err
		return model.NewAppError("CreateEmoji", "app.emojiAceess.create.internal_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil

}

func (a *App) GetPrivateEmojiImage(emojiId string, userId string) ([]byte, string, *model.AppError) {
	_, accessErr := a.Srv().Store.EmojiAccess().GetByUserIdAndEmojiId(userId, emojiId)
	if accessErr != nil {
		return nil, "", model.NewAppError("getEmojiImage", "api.emoji.get_image.read.app_error", nil, accessErr.Error(), http.StatusNotFound)
	}
	// _, storeErr := a.Srv().Store.Emoji().Get(emojiId, true)
	// if storeErr != nil {
	// 	var nfErr *store.ErrNotFound
	// 	switch {
	// 	case errors.As(storeErr, &nfErr):
	// 		return nil, "", model.NewAppError("GetEmojiImage", "app.emoji.get.no_result", nil, storeErr.Error(), http.StatusNotFound)
	// 	default:
	// 		return nil, "", model.NewAppError("GetEmojiImage", "app.emoji.get.app_error", nil, storeErr.Error(), http.StatusInternalServerError)
	// 	}
	// }

	img, appErr := a.ReadFile(getEmojiImagePath(emojiId))
	if appErr != nil {
		return nil, "", model.NewAppError("getEmojiImage", "api.emoji.get_image.read.app_error", nil, appErr.Error(), http.StatusNotFound)
	}

	_, imageType, err := image.DecodeConfig(bytes.NewReader(img))
	if err != nil {
		return nil, "", model.NewAppError("getEmojiImage", "api.emoji.get_image.decode.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return img, imageType, nil
}

func (a *App) SearchEmoji(name string, prefixOnly bool, limit int) ([]*model.Emoji, *model.AppError) {
	if !*a.Config().ServiceSettings.EnableCustomEmoji {
		return nil, model.NewAppError("SearchEmoji", "api.emoji.disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	list, err := a.Srv().Store.Emoji().Search(name, prefixOnly, limit)
	if err != nil {
		return nil, model.NewAppError("SearchEmoji", "app.emoji.get_by_name.app_error", nil, "name="+name+", "+err.Error(), http.StatusInternalServerError)
	}

	return list, nil
}

// GetEmojiStaticUrl returns a frelative static URL for system default emojis,
// and the API route for custom ones. Errors if not found or if custom and deleted.
func (a *App) GetEmojiStaticUrl(emojiName string) (string, *model.AppError) {
	subPath, _ := utils.GetSubpathFromConfig(a.Config())

	if id, found := model.GetSystemEmojiId(emojiName); found {
		return path.Join(subPath, "/static/emoji", id+".png"), nil
	}

	if emoji, err := a.Srv().Store.Emoji().GetByName(emojiName, true); err == nil {
		return path.Join(subPath, "/api/v4/emoji", emoji.Id, "image"), nil
	} else {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return "", model.NewAppError("GetEmojiStaticUrl", "app.emoji.get_by_name.no_result", nil, err.Error(), http.StatusNotFound)
		default:
			return "", model.NewAppError("GetEmojiStaticUrl", "app.emoji.get_by_name.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}
}

func resizeEmojiGif(gifImg *gif.GIF) *gif.GIF {
	// Create a new RGBA image to hold the incremental frames.
	firstFrame := gifImg.Image[0].Bounds()
	b := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())
	img := image.NewRGBA(b)

	resizedImage := image.Image(nil)
	// Resize each frame.
	for index, frame := range gifImg.Image {
		bounds := frame.Bounds()
		draw.Draw(img, bounds, frame, bounds.Min, draw.Over)
		resizedImage = resizeEmoji(img, firstFrame.Dx(), firstFrame.Dy())
		gifImg.Image[index] = imageToPaletted(resizedImage)
	}
	// Set new gif width and height
	gifImg.Config.Width = resizedImage.Bounds().Dx()
	gifImg.Config.Height = resizedImage.Bounds().Dy()
	return gifImg
}

func getEmojiImagePath(id string) string {
	return "emoji/" + id + "/image"
}

func resizeEmoji(img image.Image, width int, height int) image.Image {
	emojiWidth := float64(width)
	emojiHeight := float64(height)

	if emojiHeight <= MaxEmojiHeight && emojiWidth <= MaxEmojiWidth {
		return img
	}
	return imaging.Fit(img, MaxEmojiWidth, MaxEmojiHeight, imaging.Lanczos)
}

func imageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.Point{})
	return pm
}

func (a *App) deleteEmojiImage(id string) {
	if err := a.MoveFile(getEmojiImagePath(id), "emoji/"+id+"/image_deleted"); err != nil {
		mlog.Error("Failed to rename image when deleting emoji", mlog.String("emoji_id", id))
	}
}

func (a *App) deleteReactionsForEmoji(emojiName string) {
	if err := a.Srv().Store.Reaction().DeleteAllWithEmojiName(emojiName); err != nil {
		mlog.Warn("Unable to delete reactions when deleting emoji", mlog.String("emoji_name", emojiName), mlog.Err(err))
	}
}
