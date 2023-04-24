// Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transport

import (
	"io"
	"os"
	"regexp"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	testRootCAs = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUU1RENDQXN5Z0F3SUJBZ0lCQVRBTkJna3Fo
a2lHOXcwQkFRc0ZBREFTTVJBd0RnWURWUVFERXdkMFpYTjAKTFdOaE1CNFhEVEl6TURNeE56RXlN
ek16TTFvWERUSTBNRGt4TnpFeU5ETXpNRm93RWpFUU1BNEdBMVVFQXhNSApkR1Z6ZEMxallUQ0NB
aUl3RFFZSktvWklodmNOQVFFQkJRQURnZ0lQQURDQ0Fnb0NnZ0lCQU9oUmc1Um95UXdxCmVtQ1VN
TDU3cXVLSXJjMWZXdGdlSGRpbVRSamFsVERQMStqYUhGeG56d2M1MTRwOWNLNzcxRWZ2bDRjZW9Q
VWsKWnRhNSsxckRxdlBkd25BMnE2TXI5cFB2aWQyRkRiZVZPdXNIaHNQSG1kMDVxa1pnNGNXUGdp
eXlSM3BmNTF0bApVYkxyNU1tL0FIK0JvRHVMbFo1UG5SVUw1b0hFd1hQa3BXc0UyMXJDc2xSdmJv
WWZJYlplUzlsOHhlYURMVmdDCk53UmFHRjgxTFpoZjVrTDA0SXJUV0dETzdlbVF0S2tpN0dSZ1Ex
bHIxRHR3SXZpa0puakhBeEJiOTJ3WDN1WnoKcGdMQksxU2RsUlY1bjY2VTZtUklzMGo1MkVyTG1h
TDdUSHJxRVZHRXNvczFIbEZFQ2NJMlNhQjZZdmltaTdZawpmT1lOS2NPaE5BcXlXcWhlUERHQ0dq
d3l4RHR3OWN2Z2FJSTlTOFFUa2w5Z1JiL056dFlMREptejlEYXZiRWNjCjRDelZBdUVmdUVtWUNi
aFRrUVUyWitZczlKdXgwdmc4WXFFTExlRzlNZHc1cmZJQkkwNmRMRDVkU0JUVFc1Y08KYjRNN0h1
ODhrZUdIWnlNZXU2cVMyR2czUUFTVEM3RkpFcWFYTkRDc095aCs2Uk14UnkyZy9idEZMRm5VdmlG
QQo0NktKZHk0QWVjOEpXVkc1OFlLYkd2QlJrekkzY1BNWE1oWFpDS3pZb0tnUWoxMnFOMWM0SkVp
TUFPK2F2ZW9RCjB0dnJmd3MxMlF3d3ZIZm40SCtYVnlDZGpMcDE5dlhlY0FSRFJyaGlkRW1CbEFD
cVJVdTFLSGhzejZ2TmxzUzIKSlZiWU9BYW5ISzYzNzdYT211OUthL2x1TmxSVDdmckxBZ01CQUFH
alJUQkRNQTRHQTFVZER3RUIvd1FFQXdJQgpCakFTQmdOVkhSTUJBZjhFQ0RBR0FRSC9BZ0VBTUIw
R0ExVWREZ1FXQkJSbDk4cG1valVUQ3RZb3phTDl6L0hSCmJIUkdkREFOQmdrcWhraUc5dzBCQVFz
RkFBT0NBZ0VBNUVxL09ocGs0MUxLa3R2ajlmSWNFWXI5Vi9QUit6Z1UKQThRTUtvOVBuZSsrdW1N
dEZEK1R1M040b1lxV0srTmg0ajdHTm80RFJXejNnbWpZdUlDdklhMGo5emppT1FkWgo1Q2xVQkdk
YUlScFoyRG5CblBUM2tJbnhSd0hmU0JMVlVTRXRTcXh4YkV2dk5LWkZWY0lsRUV5ODZodnJ5OUZD
CjhFOWRXWEw5VDhMd29uVXpqSjBxZ242cGRjNHpjdEtUMDFjaDQvWGw2UjBVQkR5Q1NoSGFyU29C
eTkvSk1NTXIKajBoeEZSN3Izb052a2N3QWl6T1RsQ3BWdTZaNHF2cng3NndCc0hIanV6elNiODJL
dUxnelJUNElWbjFjbzRrVQpSaTlBRkNaRlh6QklaQlEwTUZ6NU03bzJkN0ovN3ZMOFhYRlhwWlpy
K3RibWE1L3BCSmZhcXliK3FPRXViWGdUCjFsSDZGeFNVcWt0TktQNlZoeWdQY2ZSMlR4YWtHZ0cw
Ny9qVWZWRmhpVXM5aFBlejh6Sjg2RWMrd283VEVQbEsKSlRnMHZmMDM4MTROR3ZuWmlpTnBFWVBM
S0ZhcHlDMWJONVdFTGFTWFVBaVFPZDJjK01xVHAyN21vV1RZa29TOApzRFczRTMraEN6c1djdmFY
RW1nMjZJTjQybmVUWFBuNS9QajNpcUVoT0pQYkJsY3l6dDBZL1BYeU1jR3JtbUs1CkhxOUMzTndl
VUV3M09rY09BOXlCdC9kLzZ5S3c3QmovSlFQZGI0aDlWWjNGN09wemFpeXQ5cFhvSXRQMHNUSHUK
S2ZKbDBCRUFYV29SR2lWM2EyeUlUcGp0a0pkQVBoS0xpSkkrWWowZEVEU05WZnlENFhJTXdQSmpV
eFpsd2FROQorQUtkVDFBdlplbz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`
	emptyRootCAs = ``

	testClientCert = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVORENDQWh5Z0F3SUJBZ0lSQU9JSlZ2NnB3
a0lIK0p1NTNKSEFuam93RFFZSktvWklodmNOQVFFTEJRQXcKRWpFUU1BNEdBMVVFQXhNSGRHVnpk
QzFqWVRBZUZ3MHlNekF6TVRjeE1qTTJNelphRncweU5EQTVNVGN4TWpRegpNamxhTUJFeER6QU5C
Z05WQkFNVEJtTnNhV1Z1ZERDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDCkFRb0Nn
Z0VCQU5LVFNHeEEyV2RyNmtCR0N3RjY5c1JVZElPV0xqeTUvN3QyRktKWDVVenNyMDlFWW9tS0sr
bVQKdWF2eWJIMWhsbTYzdG5kb3VFOHFIQnVhYmYvUGIzSlRTQ0twR0NRdHR2NmQzeGc3MHFZVWIx
cUZKT2o5andlNgpRZW0xb2RIVFpLc0xMc2J1N1Fzei91MGtseUovMHNYcFQ5K2JXK1M0OHMrL3pK
dHNDR21SdVhlRjE2Y1FqOWErCkFFejNqVzhrdExMYi9nS25GWGRSS2FiY2RWLzNzN2RLNWx0SXpS
ZlRvUWw0bzBpckpOa3Z4eXIrYUtMMTR4NUQKc3g2Wm9DUHJhRFYrWWlRS0ZSenFjQ1RYcWdRb3BY
LzFINFRMV3RkeG14M25IdmhZdzB0VlBZSXZsa245NmpJUwpKdVE2K1VMbVAzZDNzNWJadlhQeUZD
bENKSENxaWZNQ0F3RUFBYU9CaFRDQmdqQU9CZ05WSFE4QkFmOEVCQU1DCkE3Z3dIUVlEVlIwbEJC
WXdGQVlJS3dZQkJRVUhBd0VHQ0NzR0FRVUZCd01DTUIwR0ExVWREZ1FXQkJTRWVIOTEKVnBhdDlV
SWlrRUdkc1ljdUI2dWxOakFmQmdOVkhTTUVHREFXZ0JSbDk4cG1valVUQ3RZb3phTDl6L0hSYkhS
RwpkREFSQmdOVkhSRUVDakFJZ2daamJHbGxiblF3RFFZSktvWklodmNOQVFFTEJRQURnZ0lCQUQv
TnZVNWRSajlHCmMzYzVVQ3VLcDI4U2lacjAySE40M091WU5QVlk4L1c5QnZUSk5yMXRPRDFscnhE
eFMzTkpVdzdGaTNidmU5enMKSzA0a09peUxpVjRLd0g2eitpVm8xZU9GUzJLd1BRaGxsaDlobVBB
dXZ4Zm5Fd2k2ZEdXZm5nNExmQ1FvbXFkTgpmbkFCODJBbTViZTBubGJvaGdLcFJUWnVBZjR4dVY4
SWxlQ1pjVHdFL1hBbERhNVhHaDNvWlE3REYrQnFLSkNUCk1pYS9MT0JPYXRoRVh5ZGJmbndOUUhy
UWlQZzk4c2NMc3FTZEFQMFNGYjMrMmdscFJZT1JrQlFvOWRoa1pGZXkKc2tUakVhbk9YaUhqWldq
aXZRS2Z2WEUvK1l2eGpCcEJqREE2NnYyeUgzSlJqZEM5ZTR2cnE2R0t6VXZML3ltOQpVOGdVWnho
L2ZmeFp4TVA5UmxXajQ0U1NGUVpZNGxUNFF5U2lteFpGdVBTamwzV29QME12UHVvUzFUUzhQUk5s
CnVGeXBVell5SEtlbHpLUnRJZmlnWG9XQi9uR2hSV0RMN2FZS0xYZWRIU0ZrdXBmZm9YM1hHQThM
ZVAwQ01PaEsKUUJaUkxIeXU0VjhvRG1lakFIcFoyVjlpY2E1emtmcnJWVXFvSzF1VjYvdHd3cEZG
WDErN0w1bk0ybDJDQWxvegpaVHFUZzNCdVdYd2VkYzZQbkpuU2xQSDNadFhqcGFJUWhXdU85TUlG
WFVtVFBlSkZ2WGxKeWRsdUxtMlQzanVqCldiVENGcEhyMXBrMGk3K1J4ZVRBcFY0RTk2S09DOXEw
ZGREOG1waTM0cnkyZjFmQ2RZekhQM0s4bW5od3BPWmkKaG1Xd3VWVDV3em5kVWVBRGNWYUY2UlhU
UENKSElLd24KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`
	invalidClientCert = `bm90IGEgQ2xpZW50IENlcnRpZmljYXRlCg==`
	emptyClientCert   = ``

	testClientKey = `LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBMHBOSWJFRFpa
MnZxUUVZTEFYcjJ4RlIwZzVZdVBMbi91M1lVb2xmbFRPeXZUMFJpCmlZb3I2Wk81cS9Kc2ZXR1di
cmUyZDJpNFR5b2NHNXB0Lzg5dmNsTklJcWtZSkMyMi9wM2ZHRHZTcGhSdldvVWsKNlAyUEI3cEI2
YldoMGROa3F3c3V4dTd0Q3pQKzdTU1hJbi9TeGVsUDM1dGI1TGp5ejcvTW0yd0lhWkc1ZDRYWApw
eENQMXI0QVRQZU5ieVMwc3R2K0FxY1ZkMUVwcHR4MVgvZXp0MHJtVzBqTkY5T2hDWGlqU0tzazJT
L0hLdjVvCm92WGpIa096SHBtZ0krdG9OWDVpSkFvVkhPcHdKTmVxQkNpbGYvVWZoTXRhMTNHYkhl
Y2UrRmpEUzFVOWdpK1cKU2YzcU1oSW01RHI1UXVZL2QzZXpsdG05Yy9JVUtVSWtjS3FKOHdJREFR
QUJBb0lCQUJFSVVzSlcySDd5RHFlVwpRc3VpMjVUejA5elU1L2FIZ1BUenp5VjJnSmloU0dqYitq
QnYyYTl5QUlHMUFTdC9Ha0RvWVR6MVhuc2d4OWMvCnZZZ0VpbG92L0ZTNVlyZUNieHZYUHpWaG1W
OVBwZFlua04yN3JMY09UTWlQcFlBb1hpc3JvMlA1N1hpTGd5SkIKWkd3bzlLNkhlYXQza0k1R20z
Vk1hVXRsQ0tVcE84cUwzcEZ4S1AwMVVwbGh6ZjhMbXJpTUJQMDlxdFFJejBydQpiR1l5eUdVdk9a
a0RKZFJycmlSWGJWK0RNMFlmbVpqU1Q4aEI0UDlsOEhwMEZRNUp2TWVGREpzRFFaZjVBZnJmClFI
WE55SlFUeTNTeXJ1bGd5N0p4MGY1T2JpVWRMRWViQVRpN3VLR3Y5UEZRRUJmSzdFdE4vZ1ZibGsx
MzRzNUIKWEhkNXU1a0NnWUVBNDBVMjhONko4QXIwY2puYnNLUUJtOGhURWlJSjk3TEJPOU5kOTlJ
M1dJYklZVzIzVE5wVwo0M2R4K1JHelA4eVMzYzZhN00wbzR1dUl6TXFDSkV3cVNJUjAvVGZaWWdx
cGtwcFZPalp2VFdCUDFtSUlKUFpwCll1SFk0UVRJdkdhcVFNNnFWQXA4MW9YdXoxTmNmQWpTLzNJ
Z1BWdGVZeDNKd0pmNWVqenZQclVDZ1lFQTdUSEwKR3VCTWpqTWVhaWk1ZU1sU1BndkYxMHJISUs0
RzZCZUJDTFFXU2ViNmNOT2x2a1RaOTNqdlFiWko1L3JBTGNWNgpaTVdqbWY5Tkl0NWdDdyt2K2dM
Qm9BZXM3WEk2K2Rpdk1DYXE0dUFmWkhJWjBYbXpIOGx1a0o5ZUtyK2NyR2FzClNhWkdKRnlyQTZz
WGdOc1ZJUm85RkFsR3V1dGZnd2hSUmo1eFp3Y0NnWUVBZ241MWcyeGtDMTVlNlU5clkwdG8KV1Fo
M0dreE5LTnFNdFVzeUExL0N3NlB3WG5EZTlOUFJYQjV6WkszVEhHamNVMXVUL1MvM3NBUEpzcno4
YU5jSwoyRVNsMzljM2pHSE82QXlScnpFZVMzRm5waEwzMWpGZVpaYUVMdi9PT3M5QUpxSURqdW5P
c0dhS3JxU1F6KzlKCko3OWgzNWtjNHhCeGpaSTFmd2lKM3BrQ2dZRUFwUnBOMkExYy9IWlVxMnho
ZmRRVXJSK2d2TFZPV2s4SWU3RXcKbmhCTW0zQnR6dTlqcFVkanVVQ3l1YmpiUk9CanVQaUdzM0pt
NktDdTNxQ1BsZU43aUxrMmNlQWwzTG53bDB6ZQoxTk4xaTZxWjcxOEUzYXlxcEd1ZnpJZENFdHVC
Z1BlTzRVMGQ4ZDJYSkZ5SlphWVoxUXJnalB2UUFmZ29hWnIyCmg4Q2JTeTBDZ1lFQW1VQ3BqR0JW
MGNpVnlmUXNmOGdsclNOdWx6NzBiaVJWQzVSeno0dVJEMkhsYVM2eC8wc0IKQzltSUhpdWgwR0Zp
dEVFRlg4TzdlZ1ppNWJKMGFuQWYyakk1R1RnTjJOYzFpVlZnWldxcHh2aXpuckpKcENSYgpaejB0
M2thTkkyNjg0WTNxS2JxeG8ramRNK05hMG1qd2ErTEFOcEdCUDNwb2c0RHJ4eTNNSFdZPQotLS0t
LUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`
	invalidClientKey = `bm90IGEgQ2xpZW50IENlcnRpZmljYXRlIEtleQo=`
	emptyClientKey   = ``

	invalidBase64 = `😀`
)

var (
	missingRootCA     = regexp.MustCompile("^" + ErrRootCAMissing.Error() + "$")
	missingClientCert = regexp.MustCompile(ErrClientCertMissing.Error() + "$")
	illegalBase64Data = regexp.MustCompile(`illegal base64 data at input byte (.*)$`)
	noPEMData         = regexp.MustCompile(`failed to find any PEM data in (certificate|key) input$`)
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		config       config.Tls
		fail         bool
		errorMessage *regexp.Regexp
	}{
		{
			name: "insecure",
			config: config.Tls{
				Insecure: true,
			},
		},
		{
			name: "secure no client cert+key",
			config: config.Tls{
				Insecure: false,
				RootCas:  &testRootCAs,
			},
		},
		{
			name: "full secure",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &testRootCAs,
				ClientCert: &testClientCert,
				ClientKey:  &testClientKey,
			},
		},
		{
			name: "root CA only",
			config: config.Tls{
				Insecure: false,
				RootCas:  &testRootCAs,
			},
		},
		{
			name: "client cert+key empty",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &testRootCAs,
				ClientCert: &emptyClientCert,
				ClientKey:  &emptyClientKey,
			},
		},
		{
			name: "missing root-cas",
			config: config.Tls{
				Insecure:   false,
				ClientCert: &testClientCert,
				ClientKey:  &testClientKey,
			},
			fail:         true,
			errorMessage: missingRootCA,
		},
		{
			name: "missing client-cert",
			config: config.Tls{
				Insecure:  false,
				RootCas:   &testRootCAs,
				ClientKey: &testClientKey,
			},
			fail:         true,
			errorMessage: missingClientCert,
		},
		{
			name: "missing client-key",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &testRootCAs,
				ClientCert: &testClientCert,
			},
			fail:         true,
			errorMessage: missingClientCert,
		},
		{
			name: "empty root-cas",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &emptyRootCAs,
				ClientCert: &testClientCert,
				ClientKey:  &testClientKey,
			},
			fail:         true,
			errorMessage: missingRootCA,
		},
		{
			name: "empty client-cert",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &testRootCAs,
				ClientCert: &emptyClientCert,
				ClientKey:  &testClientKey,
			},
			fail:         true,
			errorMessage: missingClientCert,
		},
		{
			name: "empty client-key",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &testRootCAs,
				ClientCert: &testClientCert,
				ClientKey:  &emptyClientKey,
			},
			fail:         true,
			errorMessage: missingClientCert,
		},
		{
			name: "invalid root-cas base64 data",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &invalidBase64,
				ClientCert: &testClientCert,
				ClientKey:  &testClientKey,
			},
			fail:         true,
			errorMessage: illegalBase64Data,
		},
		{
			name: "invalid client-cert base64 data",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &testRootCAs,
				ClientCert: &invalidBase64,
				ClientKey:  &testClientKey,
			},
			fail:         true,
			errorMessage: illegalBase64Data,
		},
		{
			name: "invalid client-key base64 data",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &testRootCAs,
				ClientCert: &testClientCert,
				ClientKey:  &invalidBase64,
			},
			fail:         true,
			errorMessage: illegalBase64Data,
		},
		{
			name: "invalid client-cert",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &testRootCAs,
				ClientCert: &invalidClientCert,
				ClientKey:  &testClientKey,
			},
			fail:         true,
			errorMessage: noPEMData,
		},
		{
			name: "invalid client-key",
			config: config.Tls{
				Insecure:   false,
				RootCas:    &testRootCAs,
				ClientCert: &testClientCert,
				ClientKey:  &invalidClientKey,
			},
			fail:         true,
			errorMessage: noPEMData,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			transport, err := New(tc.config)
			if tc.fail {
				if assert.Error(t, err) {
					assert.Regexp(t, tc.errorMessage, err.Error())
				}
				return
			}
			assert.NoError(t, err)
			if assert.NotNil(t, transport) &&
				assert.NotNil(t, transport.TLSClientConfig) {
				assert.Equal(t, tc.config.Insecure, transport.TLSClientConfig.InsecureSkipVerify)
			}
		})
	}
}
