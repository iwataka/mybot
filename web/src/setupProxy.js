const { createProxyMiddleware } = require("http-proxy-middleware");

module.exports = function (app) {
  let target = process.env.MYBOT_BACKEND_URL || "http://localhost:8080";
  app.use(
    "/api",
    createProxyMiddleware({
      target: target,
      changeOrigin: true,
    })
  );
};
