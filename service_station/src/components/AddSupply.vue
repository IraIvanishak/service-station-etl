<style scoped>
.supply-container {
    background-color: #fff;
    border-radius: 5px;
    padding: 10px 20px;
}

.supply-header {
    color: #fff;
    background-color: #911247;
    /* Сірий фон для заголовку */
    border-radius: 5px;
    /* Закруглені кути */
    padding: 10px 20px;
    /* Внутрішній відступ */
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
    /* Відстань між заголовком та тілом */
}

.supply-body {
    background-color: #fff5f5;
    border: 1px solid #f0c0cb;
    border-radius: 10px;
    margin-bottom: 20px;
    padding: 20px;
    /* Spacing between items */
    break-inside: avoid;
    /* Prevents items from splitting across columns */
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.detail-item {
    border-bottom: 1px solid #e1e1e1;
    /* Лінія між деталями */
    padding: 10px 0;
    /* Відступи навколо деталі */
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.detail-description {
    /* Стиль для опису деталі */
    color: #333;
    /* Темний текст */
    flex: 1;
    /* Займає всю доступну ширину */
    margin-right: 20px;
    /* Відстань від ціни */
}

.detail-quantity,
.detail-price {
    /* Стилі для кількості та ціни */
    white-space: nowrap;
    /* Текст у один рядок */
}

.supplier-employee-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 10px;
    /* Відстань від деталей поставки */
}

.employee-name {
    font-style: italic;
    /* Курсив для імені співробітника */
    color: #555;
    /* Темно-сірий текст */
}

.status-tag {
    background-color: #ffffff;
    color: #2d572d;
    padding: 5px 10px;
    margin: 10px -5px;
    margin-top: 20px;
    margin-bottom: 10px;
    border-radius: 5px;
    /* Закруглені кути */
    font-size: 0.9em;
    /* Менший шрифт для статусу */
}

span {
    /* Стилі для всіх текстових елементів */
    margin-right: 10px;
    /* Відстань від сусідніх елементів */
}

.supplier-employee-section {
    flex-direction: column;
    align-items: flex-start;
}


.supply-form {
    display: flex;
    flex-direction: column;
    gap: 15px;
}

.supply-form h2 {
    color: #911247;
    margin-bottom: 10px;
}

.supply-form label {
    color: #333;
    font-size: 14px;
    margin-bottom: 5px;
}






.fas {
    margin-right: 5px;
    /* додає відстань між іконкою і текстом */
    color: #911247;
    /* встановлює колір іконок, щоб відповідати дизайну */
}

.form-container {
    background-color: #fff5f5;
    border: 1px solid #f0c0cb;
    border-radius: 10px;
    margin: 20px;
    padding: 20px;
}

.supply-form {
    display: flex;
    justify-content: space-between;
    width: 100%;
}



.form-group {
    margin-bottom: 20px;
}

.detail-list {
    width: 100%;
    padding:  0 14px;
}

.detail-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}


.supply-form-container {
    border: 1px solid #ced4da;

    margin: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px;
    background-color: #ffffff;
    border-radius: 10px;
}

.form-section {
    display: flex;
    justify-content: space-between;
    width: 100%;
    margin-bottom: 20px;
}

.form-row {
    display: flex;
    flex-direction: column;
    flex-basis: 48%;
    /* Adjust based on gap */
    margin-bottom: 10px;
}

label {
    margin-bottom: 5px;
    font-weight: bold;
    color: #333;
}

input[type="date"] {
    padding: 7px;
    border: 1px solid #ccc;
    border-radius: 4px;
    background: white;
    margin: -2px;
}

input[type="number"] {
    padding: 11px;
    border: 1px solid #ccc;
    border-radius: 4px;
    background: white;
}

select {
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    background: white;
}


.detail-item-add {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

button {
    width: 100px;
    padding: 10px 15px;
    background-color: #911247;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
    margin-left: 14px;
}

.edit-form-container {
    background-color: white;
    border: 1px solid #ced4da;
    padding: 40px;
    margin: 20px;
    margin-left: 0;

    border-radius: 20px;
}
</style>

<template>
    <div class="edit-form-container">
        <div>
            <div class="form-section">
                <div class="form-row">
                    <label for="supplier">Supplier</label>
                    <select id="supplier" v-model="newSupply.SupplierID">
                        <option disabled value="">Please select one</option>
                        <option v-for="supplier in suppliers" :key="supplier.SupplierID" :value="supplier.SupplierID">{{
                            supplier.Name }}</option>
                    </select>
                </div>
                <div class="form-row">
                    <label for="supplyDate">Supply Date</label>
                    <input type="date" id="supplyDate" v-model="newSupply.supplyDate">
                </div>
            </div>
            <div class="form-section">
                <div class="form-row">
                    <label for="manager">Manager</label>
                    <select id="manager" v-model="newSupply.EmployeeID">
                        <option disabled value="">Please select one</option>
                        <option v-for="manager in managers" :key="manager.EmployeeID" :value="manager.EmployeeID">{{
                            manager.FirstName }} {{ manager.LastName }}
                        </option>
                    </select>
                </div>
                <div class="form-row">
                    <label for="status">Status</label>
                    <select id="status" v-model="newSupply.status">
                        <option disabled value="">Please select one</option>
                        <option v-for="status in statuses" :key="status" :value="status">{{ status }}
                        </option>
                    </select>
                </div>
            </div>
            <div class="detail-list">
                <label>Detail List</label>
                <div v-for="(detail, index) in newSupply.details" :key="index" class="detail-item">
                    <select v-model="detail.DetailID">
                        <option disabled value="">Select Detail</option>
                        <option v-for="item in details" :key="item.DetailID" :value="item.DetailID">
                            {{ item.Description }}
                        </option>
                    </select>
                    <input type="number" v-model="detail.quantity" placeholder="Quantity">
                    <select v-model="detail.StationID">
                        <option disabled value="">Select Station</option>
                        <option v-for="station in stations" :key="station.StationID" :value="station.StationID">
                            {{ station.Address }}
                        </option>
                    </select>

                    <button @click="removeDetail(index)">Remove</button>
                </div>
                <div class="detail-item-add">
                    <select v-model="newDetail.DetailID">
                        <option disabled value="">Select Detail</option>
                        <option v-for="item in details" :key="item.DetailID" :value="item.DetailID">
                            {{ item.Description }}
                        </option>
                    </select>
                    <input type="number" v-model="newDetail.quantity" placeholder="Quantity">
                    <select v-model="newDetail.StationID">
                        <option disabled value="">Select Station</option>
                        <option v-for="station in stations" :key="station.StationID" :value="station.StationID">
                            {{ station.Address }}</option>
                    </select>

                    <button @click="addDetail">Add</button>
                </div>
            </div>
            <button @click="AddNewSupply" type="submit" class="submit-button">Submit</button>
        </div>
    </div>
</template>

<script>
import { API_URL } from '@/api';
export default {
    data() {
        return {
            repairs: [],
            suppliers: [],
            stations: [],
            statuses: [],
            details: [],
            managers: [],
            newDetail: {},
            newSupply: {},
            statuses: ['Pending', 'Delivered'],

        };
    },
    mounted() {
        this.init();
    },
    methods: {
        async AddNewSupply() {
            try {
                this.newSupply.supplyDate = new Date(this.newSupply.supplyDate).toISOString();
                console.log(JSON.stringify(this.newSupply))
                const response = await fetch(`${API_URL}/supplies`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(this.newSupply)
                });
            } catch (error) {
                console.error(error);
            }
            this.newSupply = { status: '', details: [] };
            this.newDetail = {
                quantity: null,
            };
        },
        addDetail() {
            this.newSupply.details.push({ ...this.newDetail });
            this.newDetail = {
                quantity: null,
            };
        },
        async init() {
            this.newSupply = { status: '', details: [] };

            try {
                const response = await fetch(`${API_URL}/suppliers`);
                const suppliers = await response.json();
                this.suppliers = suppliers;
            } catch (error) {
                console.error(error);
            }

            try {
                const response = await fetch(`${API_URL}/employees`);
                const managers = await response.json();
                this.managers = managers;
            } catch (error) {
                console.error(error);
            }

            try {
                const response = await fetch(`${API_URL}/stations`);
                const stations = await response.json();
                this.stations = stations;
            } catch (error) {
                console.error(error);
            }

            try {
                const response = await fetch(`${API_URL}/details`);
                const details = await response.json();
                this.details = details;
            } catch (error) {
                console.error(error);
            }
        },
        removeDetail(index) {
            this.newSupply.details.splice(index, 1);
        },

    }
}

</script>