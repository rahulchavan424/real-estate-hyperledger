<template>
  <div class="container">
    <el-alert type="success">
      <p>Account ID: {{ accountId }}</p>
      <p>Username: {{ userName }}</p>
      <p>Balance: ${{ balance }}</p>
    </el-alert>
    <div v-if="sellingList.length == 0" style="text-align: center;">
      <el-alert title="No data found" type="warning" />
    </div>
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val, index) in sellingList" :key="index" :span="6" :offset="1">
        <el-card class="all-card">
          <div slot="header" class="clearfix">
            <span>{{ val.sellingStatus }}</span>
            <el-button v-if="roles[0] !== 'admin' && (val.seller === accountId || val.buyer === accountId) && val.sellingStatus !== 'Completed' && val.sellingStatus !== 'Expired' && val.sellingStatus !== 'Cancelled'" style="float: right; padding: 3px 0" type="text" @click="updateSelling(val, 'cancelled')">Cancel</el-button>
            <el-button v-if="roles[0] !== 'admin' && val.seller === accountId && val.sellingStatus === 'In Progress'" style="float: right; padding: 3px 8px" type="text" @click="updateSelling(val, 'done')">Confirm Payment</el-button>
            <el-button v-if="roles[0] !== 'admin' && val.sellingStatus === 'On Sale' && val.seller !== accountId" style="float: right; padding: 3px 0" type="text" @click="createSellingByBuy(val)">Buy</el-button>
          </div>
          <div class="item">
            <el-tag>Real Estate ID: </el-tag>
            <span>{{ val.objectOfSale }}</span>
          </div>
          <div class="item">
            <el-tag type="success">Seller ID: </el-tag>
            <span>{{ val.seller }}</span>
          </div>
          <div class="item">
            <el-tag type="danger">Price: </el-tag>
            <span>$ {{ val.price }}</span>
          </div>
          <div class="item">
            <el-tag type="warning">Validity Period: </el-tag>
            <span>{{ val.salePeriod }} days</span>
          </div>
          <div class="item">
            <el-tag type="info">Creation Time: </el-tag>
            <span>{{ val.createTime }}</span>
          </div>
          <div class="item">
            <el-tag>Buyer ID: </el-tag>
            <span v-if="val.buyer === ''">Vacant</span>
            <span>{{ val.buyer }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { querySellingList, createSellingByBuy, updateSelling } from '@/api/selling';

export default {
  name: 'AllSelling',
  data() {
    return {
      loading: true,
      sellingList: [],
    };
  },
  computed: {
    ...mapGetters([
      'accountId',
      'roles',
      'userName',
      'balance',
    ]),
  },
  created() {
    querySellingList().then(response => {
      if (response !== null) {
        this.sellingList = response;
      }
      this.loading = false;
    }).catch(_ => {
      this.loading = false;
    });
  },
  methods: {
    createSellingByBuy(item) {
      this.$confirm('Buy now?', 'Confirmation', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        type: 'success',
      }).then(() => {
        this.loading = true;
        createSellingByBuy({
          buyer: this.accountId,
          objectOfSale: item.objectOfSale,
          seller: item.seller,
        }).then(response => {
          this.loading = false;
          if (response !== null) {
            this.$message({
              type: 'success',
              message: 'Purchase successful!',
            });
          } else {
            this.$message({
              type: 'error',
              message: 'Purchase failed!',
            });
          }
          setTimeout(() => {
            window.location.reload();
          }, 1000);
        }).catch(_ => {
          this.loading = false;
        });
      }).catch(() => {
        this.loading = false;
        this.$message({
          type: 'info',
          message: 'Purchase canceled',
        });
      });
    },
    updateSelling(item, type) {
      let tip = '';
      if (type === 'done') {
        tip = 'Confirm Payment';
      } else {
        tip = 'Cancel Operation';
      }
      this.$confirm('Do you want to ' + tip + '?', 'Confirmation', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        type: 'success',
      }).then(() => {
        this.loading = true;
        updateSelling({
          buyer: item.buyer,
          objectOfSale: item.objectOfSale,
          seller: item.seller,
          status: type,
        }).then(response => {
          this.loading = false;
          if (response !== null) {
            this.$message({
              type: 'success',
              message: tip + ' successful!',
            });
          } else {
            this.$message({
              type: 'error',
              message: tip + ' failed!',
            });
          }
          setTimeout(() => {
            window.location.reload();
          }, 1000);
        }).catch(_ => {
          this.loading = false;
        });
      }).catch(() => {
        this.loading = false;
        this.$message({
          type: 'info',
          message: 'Operation ' + tip + ' canceled',
        });
      });
    },
  },
};
</script>

<style>
  .container {
    width: 100%;
    text-align: center;
    min-height: 100%;
    overflow: hidden;
  }
  .tag {
    float: left;
  }

  .item {
    font-size: 14px;
    margin-bottom: 18px;
    color: #999;
  }

  .clearfix:before,
  .clearfix:after {
    display: table;
  }
  .clearfix:after {
    clear: both;
  }

  .all-card {
    width: 280px;
    height: 380px;
    margin: 18px;
  }
</style>
