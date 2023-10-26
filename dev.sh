#!/bin/bash

artelamint_mod=./go.mod
artelasdk_mod=../aspect-core/go.mod
evm_mod=../artela-evm/go.mod
artela_aspect_mod=../artela-cosmos-sdk/go.mod

# add replacement to artelamint
update_artelamint_mod() {
    cat >>$artelamint_mod <<EOF

replace (
    github.com/artela-network/aspect-core => ../aspect-core
    github.com/artela-network/artela-evm => ../artela-evm
    github.com/artela-network/aspect-runtime => ../aspect-runtime
)
EOF
    sed -i -e "s/\(github.com\/cosmos\/cosmos-sdk => \).*/\1..\/artela-cosmos-sdk/" $artelamint_mod
    rm -f $artelamint_mod-e
}

# add replacement to artelasdk
update_artelasdk_mod() {
    cat >>$artelasdk_mod <<EOF

replace (
    github.com/artela-network/aspect-runtime => ../aspect-runtime
)
EOF
}

# add replacement to artela-aspect
update_artela_aspect_mod() {
    cat >>$artela_aspect_mod <<EOF

replace (
    github.com/artela-network/aspect-core => ../aspect-core
    github.com/artela-network/aspect-runtime => ../aspect-runtime
)
EOF
    sed -i -e "s/\(github.com\/tendermint\/tendermint => \).*/\1..\/cometbft/" $artela_aspect_mod
    rm -f $artelamint_mod-e
}


# add replacement to evm
update_evm_mod() {
    cat >>$evm_mod <<EOF

replace (
    github.com/artela-network/aspect-core => ../aspect-core
    github.com/artela-network/aspect-runtime => ../aspect-runtime
)
EOF
}

input=$1
if [ $input == "set" ]; then
    echo "setting submodule ..."
    update_artelamint_mod
    update_artelasdk_mod
    update_artela_aspect_mod
    update_evm_mod
    echo "done"
elif [ $input == "reset" ]; then
    echo "resetting go.mod ..."
    curdir=$(pwd)
    git checkout go.mod

    cd ../aspect-core
    git checkout go.mod

    cd ../artela-cosmos-sdk
    git checkout go.mod

    cd ../artela-evm
    git checkout go.mod

    echo "done"
else
    echo "nothing has been changes, dev set|reset"
fi
