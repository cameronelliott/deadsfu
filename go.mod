module github.com/x186k/sfu1

go 1.16

require (
	github.com/caddyserver/certmagic v0.12.0
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.2.0 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/libdns/cloudflare v0.0.0-20200528144945-97886e7873b1
	github.com/libdns/duckdns v0.1.0
	github.com/mholt/acmez v0.1.2 // indirect
	github.com/miekg/dns v1.1.40
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pion/dtls/v2 v2.0.8 // indirect
	github.com/pion/rtcp v1.2.6
	github.com/pion/rtp v1.6.2
	github.com/pion/sdp/v3 v3.0.4
	github.com/pion/webrtc/v3 v3.0.12
	github.com/pkg/profile v1.5.0
	github.com/stretchr/testify v1.7.0
	github.com/x186k/ddns5libdns v0.0.0-20210226014652-bb305912e632
	github.com/x186k/dynamicdns v0.0.0-20210209094553-08838551f0e5
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073 // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/tools v0.0.0-20200513154647-78b527d18275 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

//replace github.com/pion/webrtc/v3 => ../webrtc
//replace github.com/pion/webrtc/v3 v3.0.4 => github.com/cameronelliott/webrtc/v3 v3.0.5
