package translation

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	translatev3 "cloud.google.com/go/translate/apiv3"
	"golang.org/x/text/language"
	translatepbv3 "google.golang.org/genproto/googleapis/cloud/translate/v3"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func newSys(options Options) (sys *Translation, err error) {
	sys = &Translation{options: options}
	return
}

type Translation struct {
	options Options
}

//文字翻译
func (this *Translation) Translation_Text_Base(ctx context.Context, original string, from, to string) (result string, err error) {
	var (
		client *translate.Client
	)
	client, err = translate.NewClient(ctx)
	if err != nil {
		return
	}
	target, err := language.Parse(to)
	if err != nil {
		return
	}
	translations, err := client.Translate(ctx, []string{original}, target, nil)
	if err != nil {
		return
	}
	result = translations[0].Text
	return
}

//文字翻译
func (this *Translation) Translation_Text_GoogleV3(ctx context.Context, original string, from, to string) (result string, err error) {
	client, err := translatev3.NewTranslationClient(ctx)
	if err != nil {
		return
	}
	defer client.Close()
	req := &translatepbv3.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", this.options.AppId),
		SourceLanguageCode: from,
		TargetLanguageCode: to,
		MimeType:           "text/plain", //Mime types: "text/plain", "text/html"
		Contents:           []string{original},
	}
	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return
	}
	result = resp.GetTranslations()[0].GetTranslatedText()
	return
}

//语音转文字
func (this *Translation) Translation_Voice(ctx context.Context, original []byte, from string) (result string, err error) {
	var (
		client *speech.Client
		resp   *speechpb.RecognizeResponse
	)
	if client, err = speech.NewClient(ctx); err != nil {
		return
	}
	// Detects speech in the audio file.
	if resp, err = client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    from,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: original},
		},
	}); err != nil {
		return
	}
	for _, _result := range resp.Results {
		for _, alt := range _result.Alternatives {
			result += alt.Transcript
		}
	}
	return
}
