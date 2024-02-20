FROM simonalphafang/artela-client:0.0.1

RUN git clone -b feat-ci https://github.com/luketheart/test-scripts
RUN /bin/bash -c "source ~/.nvm/nvm.sh && cd /test-scripts/pre_public_testnet/json-rpc && npm i"