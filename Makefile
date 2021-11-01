# Author: ccwutw
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ==============================================================================


GOROOT ?= $(shell go env GOROOT)

DIST      := dist
OBJS      := $(addprefix $(DIST)/,$(notdir $(wildcard src/*)))
PB_SRCS   := message/*.proto
PB_OBJS   := message/*.pb.go

CERT_ROOT := $(DIST)/certificates
CSR_CONF  := csr.conf
SRLS      := $(CERT_ROOT)/ca.srl
CRTS      := $(patsubst %.crt,$(CERT_ROOT)/%.crt,ca.crt client.crt server.crt proxy.crt)
KEYS      := $(CRTS:%.crt=%.key)
PFXS      := $(CRTS:%.crt=%.pfx)

.PRECIOUS: %.key

.DEFAULT_GOAL = all

all: build

build: certificates $(PB_OBJS) $(OBJS)

certificates: $(CRTS) $(PFXS)

test: $(PB_OBJS) $(PB_SRCS)
	go test ./...

clean:
	-rm $(PB_OBJS) $(OBJS)

distclean: clean
	-rm $(PFXS) $(CRTS) $(CSRS) $(KEYS) $(SRLS)
	-rm -rf $(DIST)

# Rules

$(CERT_ROOT):
	mkdir -p $(CERT_ROOT)

$(DIST):
	mkdir -p $(DIST)

$(DIST)/%: ./src/% src/%/*.go | $(DIST)
	go build -o $@ ./$<

$(PB_OBJS): $(PB_SRCS)
	protoc -I message/ $^ --go_out=message --go_opt=paths=source_relative --go-grpc_out=message --go-grpc_opt=paths=source_relative

%.key: | $(CERT_ROOT)
	openssl ecparam -out $@ -name prime256v1 -genkey

%.csr: %.key | $(CERT_ROOT)
	openssl req -new -sha256 -key $^ -out $@ -subj '/CN=localhost'

%.crt: %.csr %.key $(CERT_ROOT)/ca.key $(CERT_ROOT)/ca.crt | $(CERT_ROOT)
	openssl x509 -req -in $< -out $@ -sha256 -days 365 -CA $(CERT_ROOT)/ca.crt -CAkey $(CERT_ROOT)/ca.key -CAcreateserial -extfile $(CSR_CONF)
	openssl x509 -in $@ -text -noout

%/ca.crt: $(CERT_ROOT)/ca.key | $(CERT_ROOT)
	openssl req -x509 -new -nodes -key $^ -sha256 -days 365 -out $@ -subj '/CN=localhost'
	openssl x509 -in $@ -text -noout

%.pfx: | $(CERT_ROOT)
	openssl pkcs12 -inkey $(@:%.pfx=%.key) -in $(@:%.pfx=%.crt) -export -nodes -passout pass: -out $@
