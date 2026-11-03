<template>
  <div class="data-load-history-container">
    <h1>Data Load Manegment</h1>
    <div class="panel">
      <div class="initial-load-group">
        <button @click="clearStorage" class="button clear-button"><i class="fas fa-trash-alt"></i> Clear
          Storage</button>
        <button @click="performInitialLoad" class="button load-button" :disabled="isData"><i
            class="fas fa-download"></i> Initial Load</button>

      </div>
      <div class="incremental-load-group">
        <button @click="fetchStage" class="button clear-button"><i class="fas fa-sync-alt"></i> Detect Changes</button>
        <button @click="performIncrementalLoad" class="button load-button" :disabled="!isData"><i
            class="fas fa-plus-circle"></i> Incremental Load</button>
      </div>

      <div v-if="stage != null" class="summary-container">
        <h3 v-if="stage.some(a => a.RowCount > 0)">
          Found changes in OLTP</h3>
        <h3 v-else>
          No changes found in OLTP</h3>

        <div class="table-summary" v-for="table in stage" :key="table.TableName">
          <div v-if="table.RowCount != 0">
            <div class="table-name">{{ table.TableName }}</div>
            <div class="row-count">{{ table.RowCount }}</div>
          </div>
        </div>
      </div>
    </div>



    <table class="data-table">
      <thead>
        <tr>
          <th>Load Datetime</th>
          <th>Load Time, ms</th>
          <th>Load Rows</th>
          <th>Affected Table Count</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in dataLoadHistory" :key="item.DataLoadHistoryID">
          <td>{{ item.LoadDatetime.Valid ? new Date(item.LoadDatetime.Time).toDateString() : 'N/A' }}</td>
          <td>{{ item.LoadTime.Valid ? new Date(item.LoadDatetime.Time).getMilliseconds() : 'N/A' }}</td>
          <td>{{ item.LoadRows.Valid ? item.LoadRows.Int32 : 'N/A' }}</td>
          <td>{{ item.AffectedTableCount.Valid ? item.AffectedTableCount.Int32 : 'N/A' }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>

h1 {
    color: #333;
}

.data-load-history-container {
  max-width: 1200px;
}

.initial-load-group,
.incremental-load-group {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-start;
  /* Aligns items to the start of the line */
}

.panel {
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


.button {
  width: 200px;
  padding: 10px 20px;
  margin-right: 10px;
  color: white;
  border-radius: 5px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
}

.button i {
  margin-right: 5px;
}

.load-button {
  background-color: #911247;
}

.load-button:hover {
  background-color: #a83250;
}

.load-button:disabled {
  background-color: #ccc;
  color: #666;
  cursor: not-allowed;
}

.clear-button {
  background-color: transparent;
  color: #911247;
  border: 2px solid #911247;
}

.clear-button:hover {
  background-color: #f5f5f5;
}

.summary-container {
  display: flex;
  flex-direction: column;
}

.table-summary {
  display: flex;
  flex-wrap: wrap;
  /* Allows items to wrap onto the next line */
  gap: 20px 10px;
  /* Vertical and horizontal spacing */
}

.table-name,
.row-count {
  width: 200px;
  /* Fixed width for consistent alignment */
  padding: 5px;
  margin-bottom: 5px;
  text-align: left;
  /* Text alignment for table name */
  background-color: #ffffff;
  /* Background color for visibility */
  border-radius: 5px;
  /* Rounded corners for aesthetics */
  display: inline-flex;
  /* Makes each name and count a flexible inline element */
  justify-content: space-between;
  /* Spreads the name and count across the available space */
}


.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 10px;
  border: 1px solid #ddd;
}

.data-table th {
  background-color: #911247;
  color: white;
}

.data-table td {
  background-color: #fff;
  color: #333;
}

.data-table tr:nth-child(odd) td {
  background-color: #f9f9f9;
}

button {
  padding: 10px 20px;
  color: white;
  background-color: #911247;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #a83250;
}
</style>


<script>
import { API_URL } from '@/api';
export default {
  data() {
    return {
      dataLoadHistory: [],
      stage: null,
      isData: true
    };
  },
  mounted() {
    this.fetchDataLoadHistory();
    this.isDataf();
  },
  computed: {
    canPerformIncrementalLoad() {
      return this.dataLoadHistory !== null && this.dataLoadHistory.length > 0;
    }
  },
  methods: {
    async isDataf() {
      try {
        const response = await fetch(`${API_URL}/supply-cube`);
        const data = await response.json();
        if (data.length > 0) {
          this.isData = true;
        } else {
          this.isData = false;
        }
      } catch (error) {
        console.error('Failed to fetch data:', error);
      }
    },
    async fetchStage() {
      try {
        const response = await fetch(`${API_URL}/stage`);

        this.stage = await response.json();
      } catch (error) {
        console.error('Error performing stage:', error);
        alert('Error performing stage.');
      }
    },
    async clearStorage() {
      try {
        const response = await fetch(`${API_URL}/clearmeta`, {
          method: 'POST' // Змініть на потрібний метод та URL
        });
        this.isData = false;
      } catch (error) {
        console.error('Error clearing storage and metadata:', error);
        alert('Error clearing storage and metadata.');
      }

    },
    async fetchDataLoadHistory() {
      try {
        const response = await fetch(`${API_URL}/meta`);
        const data = await response.json();
        this.dataLoadHistory = data;
      } catch (error) {
        console.error('Failed to fetch data load history:', error);
      }
    },
    async performInitialLoad() {
      try {
        const response = await fetch(`${API_URL}/init`, {
          method: 'POST'
        });
        if (response.ok) {
          alert('Initial load successful!');
          this.fetchDataLoadHistory(); // Оновлюємо історію завантаження
        } else {
          alert('Failed to perform initial load.');
        }
      } catch (error) {
        console.error('Error performing initial load:', error);
        alert('Error performing initial load.');
      }
      this.isDataf();
    },

    async performIncrementalLoad() {
      if (!this.canPerformIncrementalLoad) {
        alert('Incremental load is not possible without previous data.');
        return;
      }
      try {
        const response = await fetch(`${API_URL}/incremental`, {
          method: 'POST'
        });
        if (response.ok) {
          alert('Incremental load successful!');
          this.fetchDataLoadHistory(); // Оновлюємо історію завантаження
        } else {
          alert('Failed to perform incremental load.');
        }
      } catch (error) {
        console.error('Error performing incremental load:', error);
        alert('Error performing incremental load.');
      }
    }
  }
}
</script>
