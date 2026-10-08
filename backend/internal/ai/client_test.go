package ai

import "testing"

func TestAllowedHost(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{name: "allowed host", url: "https://api.yandexcloud.net/v1/completion", want: true},
		{name: "allowed explicit TLS port", url: "https://gigachat.devices.sberbank.ru:443/api", want: true},
		{name: "untrusted suffix", url: "https://api.yandexcloud.net.attacker.example/v1", want: false},
		{name: "untrusted user info", url: "https://api.yandexcloud.net@attacker.example/v1", want: false},
		{name: "insecure scheme", url: "http://api.yandexcloud.net/v1", want: false},
		{name: "non TLS port", url: "https://api.yandexcloud.net:8443/v1", want: false},
		{name: "invalid URL", url: "://api.yandexcloud.net", want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := allowedHost(test.url); got != test.want {
				t.Fatalf("allowedHost(%q) = %t, want %t", test.url, got, test.want)
			}
		})
	}
}
