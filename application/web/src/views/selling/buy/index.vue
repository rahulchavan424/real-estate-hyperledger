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
        <el-card class="buy-card">
          <div slot="header" class="clearfix">
            <span>{{ val.selling.sellingStatus }}</span>
            <el-button v-if="val.selling.sellingStatus !== 'Completed' && val.selling.sellingStatus !== 'Expired' && val.selling.sellingStatus !== 'Cancelled'" style="float: right; padding: 3px 0" type="text" @click="updateSelling(val, 'cancelled')">Cancel</el-button>
          </div>
          <div class="item">
            <el-tag type="warning">Order Time: </el-tag>
            <span>{{ val.createTime }}</span>
          </div>
          <div class="item">
            <el-tag>Real Estate ID: </el-tag>
            <span>{{ val.selling.objectOfSale }}</span>
          </div>
          <div class="item">
            <el-tag type="success">Seller ID: </el-tag>
            <span>{{ val.selling.seller }}</span>
          </div>
          <div class="item">
            <el-tag type="danger">Price: </el-tag>
            <span>$ {{ val.selling.price }} </span>
          </div>
          <div class="item">
            <el-tag type="warning">Validity Period: </el-tag>
            <span>{{ val.selling.salePeriod }} days</span>
          </div>
          <div class="item">
            <el-tag type="info">Creation Time: </el-tag>
            <span>{{ val.selling.createTime }}</span>
          </div>
          <div class="item">
            <el-tag>Buyer ID: </el-tag>
            <span v-if="val.selling.buyer === ''">Vacant</span>
            <span>{{ val.selling.buyer }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { querySellingListByBuyer, updateSelling } from '@/api/selling';

export default {
  name: 'BuySelling',
  data() {
    return {
      loading: true,
      sellingList: [],
    };
  },
  computed: {
    ...mapGetters([
      'accountId',
      'userName',
      'balance',
    ]),
  },
  created() {
    querySellingListByBuyer({ buyer: this.accountId }).then(response => {
      if (response !== null) {
        this.sellingList = response;
      }
      this.loading = false;
    }).catch(_ => {
      this.loading = false;
    });
  },
  methods: {
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
          buyer: item.selling.buyer,
          objectOfSale: item.selling.objectOfSale,
          seller: item.selling.seller,
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

  .buy-card {
    width: 280px;
    height: 430px;
    margin: 18px;
  }
</style>
