import api from "./client";

export const partnerLogin = (email, password) =>
  api.post("/auth/partner/login", { email, password }).then((r) => r.data.data);

export const getPartnerProfile = () =>
  api.get("/partner/profile").then((r) => r.data.data);

export const getPartnerDashboard = () =>
  api.get("/partner/dashboard").then((r) => r.data.data);

export const getPartnerStats = (params = {}) =>
  api.get("/partner/stats", { params }).then((r) => r.data.data);

export const changePassword = (current_password, new_password) =>
  api
    .post("/partner/change-password", { current_password, new_password })
    .then((r) => r.data.data);

export const savePartnerLegalInfo = (data) =>
  api.post("/partner/legal-info", data).then((r) => r.data.data);

export const createLocation = (data) =>
  api.post("/partner/locations", data).then((r) => r.data.data);

export const listLocations = () =>
  api.get("/partner/locations").then((r) => r.data.data);

export const listEmployees = () =>
  api.get("/partner/employees").then((r) => r.data.data);

export const createEmployee = (data) =>
  api.post("/partner/employees", data).then((r) => r.data.data);

export const deleteEmployee = (id) =>
  api.delete(`/partner/employees/${id}`).then((r) => r.data.data);

// Boxes
export const createBox = (data) =>
  api.post("/partner/boxes", data).then((r) => r.data.data);

export const getBoxById = (id) =>
  api.get(`/partner/boxes/${id}`).then((r) => r.data.data);

export const updateBox = (id, data) =>
  api.patch(`/partner/boxes/${id}`, data).then((r) => r.data.data);

export const deleteBox = (id) =>
  api.delete(`/partner/boxes/${id}`).then((r) => r.data.data);

export const listBoxes = () =>
  api.get("/partner/boxes").then((r) => r.data.data);

export const listPendingConfirmationOrders = () =>
  api.get("/partner/orders/pending-confirmation").then((r) => r.data.data);

export const listLocationBoxes = (locationId) =>
  api.get(`/locations/${locationId}/boxes`).then((r) => r.data.data);

// Orders
export const listPartnerOrders = ({ status = "", limit = 20, offset = 0 } = {}) =>
  api
    .get("/partner/orders", {
      params: { status, limit, offset },
    })
    .then((r) => r.data.data);

export const getOrderByPickupCode = (pickupCode) =>
  api
    .get(`/partner/orders/by-code/${encodeURIComponent(pickupCode)}`)
    .then((r) => r.data.data);

export const pickupOrder = (orderId) =>
  api.post(`/partner/orders/${orderId}/pickup`).then((r) => r.data.data);

// Media
export const uploadMedia = (file) => {
  const formData = new FormData();
  formData.append("file", file);
  return api
    .post("/media/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    })
    .then((r) => r.data.data);
};
