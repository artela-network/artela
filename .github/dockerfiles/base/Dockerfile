FROM ubuntu
WORKDIR /
RUN apt update && apt -y upgrade && apt -y install build-essential cmake libboost-all-dev curl git
RUN git clone https://github.com/artela-network/artela-web3.js.git

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash 

RUN /bin/bash -c "source ~/.nvm/nvm.sh && nvm install v16.17.0"
RUN /bin/bash -c "source ~/.nvm/nvm.sh && cd /artela-web3.js && npm i && npm run build"

RUN git clone https://github.com/artela-network/artela-solidity.git && sed -i '705s/if (Expression const\* value = _return.expression())/if (_return.expression())/g' /artela-solidity/libsolidity/codegen/ir/IRGeneratorForStatements.cpp
RUN cd /artela-solidity && cmake -B build . && make -j4 -C build && make -C build install

RUN SOLC_PATH=$(which solc) && ASOLC_PATH=$(echo $SOLC_PATH | sed "s/solc/asolc/g") && ln -s $SOLC_PATH $ASOLC_PATH



# RUN git clone https://github.com/artela-network/aspect-tooling.git
# RUN /bin/bash -c "source ~/.nvm/nvm.sh && cd /aspect-tooling/packages/testcases && npm i && npm run contract:build"
# RUN /bin/bash -c "source ~/.nvm/nvm.sh && cd /aspect-tooling/packages/testcases && sh build.sh || true"