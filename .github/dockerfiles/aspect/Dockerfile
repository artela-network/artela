FROM simonalphafang/artela-client:0.0.1

RUN git clone https://github.com/artela-network/aspect-tooling.git
RUN /bin/bash -c "source ~/.nvm/nvm.sh && cd /aspect-tooling/packages/testcases && npm i && npm run contract:build"
RUN /bin/bash -c "source ~/.nvm/nvm.sh && cd /aspect-tooling/packages/testcases && sh build.sh || true"