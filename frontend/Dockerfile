FROM node:22-alpine AS builder

WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .

ARG VITE_API_BASE_URL
ARG VITE_CHAIN_ID
ARG VITE_DID_PROFILE_ADDRESS
ARG VITE_EVIDENCE_REGISTRY_ADDRESS
ARG VITE_CREDENTIAL_SBT_ADDRESS
ARG VITE_REPUTATION_ADDRESS

RUN npm run build

FROM nginx:1.25-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
