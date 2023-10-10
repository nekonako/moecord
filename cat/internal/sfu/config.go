package sfu

import (
	"github.com/jiyeyuran/mediasoup-go"
	"github.com/jiyeyuran/mediasoup-go/h264"
)

var consumerDeviceCapabilities = mediasoup.RtpCapabilities{
	Codecs: []*mediasoup.RtpCodecCapability{
		{
			MimeType:             "audio/opus",
			Kind:                 "audio",
			PreferredPayloadType: 100,
			ClockRate:            48000,
			Channels:             2,
		},
		{
			MimeType:             "video/H264",
			Kind:                 "video",
			PreferredPayloadType: 101,
			ClockRate:            90000,
			Parameters: mediasoup.RtpCodecSpecificParameters{
				RtpParameter: h264.RtpParameter{
					LevelAsymmetryAllowed: 1,
					PacketizationMode:     1,
					ProfileLevelId:        "4d0032",
				},
			},
			RtcpFeedback: []mediasoup.RtcpFeedback{
				{Type: "nack", Parameter: ""},
				{Type: "nack", Parameter: "pli"},
				{Type: "ccm", Parameter: "fir"},
				{Type: "goog-remb", Parameter: ""},
			},
		},
		{
			MimeType:             "video/rtx",
			Kind:                 "video",
			PreferredPayloadType: 102,
			ClockRate:            90000,
			Parameters: mediasoup.RtpCodecSpecificParameters{
				Apt: 101,
			},
		},
	},
}

var mediasoupRouterConfig = mediasoup.RouterOptions{
	MediaCodecs: []*mediasoup.RtpCodecCapability{
		{
			Kind:      mediasoup.MediaKind_Audio,
			MimeType:  "audio/opus",
			ClockRate: 48000,
			Channels:  2,
		},
		{
			Kind:      mediasoup.MediaKind_Video,
			MimeType:  "video/VP8",
			ClockRate: 90000,
			Parameters: mediasoup.RtpCodecSpecificParameters{
				XGoogleStartBitrate: 1000,
			},
		},
		{
			Kind:      mediasoup.MediaKind_Video,
			MimeType:  "video/VP9",
			ClockRate: 90000,
			Parameters: mediasoup.RtpCodecSpecificParameters{
				ProfileId:           "2",
				XGoogleStartBitrate: 1000,
			},
		},
		{
			Kind:      mediasoup.MediaKind_Video,
			MimeType:  "video/h264",
			ClockRate: 90000,
			Parameters: mediasoup.RtpCodecSpecificParameters{
				RtpParameter: h264.RtpParameter{
					PacketizationMode:     1,
					ProfileLevelId:        "42e01f",
					LevelAsymmetryAllowed: 1,
				},
				XGoogleStartBitrate: 1000,
			},
		},
		{
			Kind:      mediasoup.MediaKind_Video,
			MimeType:  "video/h264",
			ClockRate: 90000,
			Parameters: mediasoup.RtpCodecSpecificParameters{
				RtpParameter: h264.RtpParameter{
					PacketizationMode:     1,
					ProfileLevelId:        "4d0032",
					LevelAsymmetryAllowed: 1,
				},
				XGoogleStartBitrate: 1000,
			},
		},
	},
}
