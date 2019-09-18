# Copyright 1999-2019 Gentoo Authors
# Distributed under the terms of the GNU General Public License v2

EAPI=7

DESCRIPTION="Command line tool to bootstrap Github repository"
HOMEPAGE="https://github.com/g4s8/${PN}"
SRC_URI="${HOMEPAGE}/archive/${PV}.tar.gz"

LICENSE="MIT"
SLOT="0"
KEYWORDS="~amd64 ~x86"
IUSE="test"

DEPEND="dev-lang/go"
RDEPEND="dev-vcs/git"

_go_get() {
	local name=$1
	go get -v -u $name || die "'go get' failed to get $name"
}

src_prepare() {
	default
	mkdir -pv $HOME/go/src/github.com/g4s8 || die
	ln -snv $PWD $HOME/go/src/github.com/g4s8/gitstrap || die
	_go_get github.com/google/go-github/github
	_go_get golang.org/x/oauth2
	_go_get gopkg.in/yaml.v2
}

src_compile() {
	local now=$(date -u +%Y.%m.%dT%H:%M:%S)
	emake OUTPUT=lib BUILD_VERSION=${PV} BUILD_DATE=${now} lib || die
	go build -o ${PN} ./cmd/${PN} || die
}

src_test() {
	emake OUTPUT=lib test || die
}

src_install() {
	dobin ${PN} || die "${PN} installation failed"
	elog "Read the README for details: ${HOMEPAGE}"
}
