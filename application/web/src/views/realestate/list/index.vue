<template>
  <div class="container">
    <el-alert type="success">
      <p>Account ID: {{ accountId }}</p>
      <p>Username: {{ userName }}</p>
      <p>Balance: ${{ balance }}</p>
      <p>When initiating a sale, donation, or pledge, the collateral status is true</p>
      <p>You can initiate a sale, donation, or pledge operation only when the collateral status is false</p>
    </el-alert>
    <div v-if="realEstateList.length == 0" style="text-align: center;">
      <el-alert title="No data found" type="warning" />
    </div>
    <el-row v-loading="loading" :gutter="20">
      <el-col v-for="(val, index) in realEstateList" :key="index" :span="6" :offset="1">
        <el-card class="realEstate-card">
          <div slot="header" class="clearfix">
            Collateral Status:
            <span style="color: rgb(255, 0, 0);">{{ val.encumbrance }}</span>
          </div>

          <div class="item">
            <el-tag>Real Estate ID: </el-tag>
            <span>{{ val.realEstateId }}</span>
          </div>
          <div class="item">
            <el-tag type="success">Owner ID: </el-tag>
            <span>{{ val.proprietor }}</span>
          </div>
          <div class="item">
            <el-tag type="warning">Total Area: </el-tag>
            <span>{{ val.totalArea }} m²</span>
          </div>
          <div class="item">
            <el-tag type="danger">Living Space: </el-tag>
            <span>{{ val.livingSpace }} m²</span>
          </div>

          <div v-if="!val.encumbrance && roles[0] !== 'admin'">
            <el-button type="text" @click="openDialog(val)">Sell</el-button>
            <el-divider direction="vertical" />
            <el-button type="text" @click="openDonatingDialog(val)">Donate</el-button>
          </div>
          <el-rate v-if="roles[0] === 'admin" />
        </el-card>
      </el-col>
    </el-row>
    <el-dialog v-loading="loadingDialog" :visible.sync="dialogCreateSelling" :close-on-click-modal="false" @close="resetForm('realForm')">
      <el-form ref="realForm" :model="realForm" :rules="rules" label-width="100px">
        <el-form-item label="Price (USD)" prop="price">
          <el-input-number v-model="realForm.price" :precision="2" :step="10000" :min="0" />
        </el-form-item>
        <el-form-item label="Validity Period (days)" prop="salePeriod">
          <el-input-number v-model="realForm.salePeriod" :min="1" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="createSelling('realForm')">Sell Now</el-button>
        <el-button @click="dialogCreateSelling = false">Cancel</el-button>
      </div>
    </el-dialog>
    <el-dialog v-loading="loadingDialog" :visible.sync="dialogCreateDonating" :close-on-click-modal="false" @close="resetForm('DonatingForm')">
      <el-form ref="DonatingForm" :model="DonatingForm" :rules="rulesDonating" label-width="100px">
        <el-form-item label="Owner" prop="proprietor">
          <el-select v-model="DonatingForm.proprietor" placeholder="Select an owner" @change="selectGet">
            <el-option v-for="item in accountList" :key="item.accountId" :label="item.userName" :value="item.accountId">
              <span style="float: left">{{ item.userName }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">{{ item.accountId }}</span>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="createDonating('DonatingForm')">Donate Now</el-button>
        <el-button @click="dialogCreateDonating = false">Cancel</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { queryAccountList } from '@/api/account';
import { queryRealEstateList } from '@/api/realEstate';
import { createSelling } from '@/api/selling';
import { createDonating } from '@/api/donating';

export default {
  name: 'RealEstate',
  data() {
    var checkArea = (rule, value, callback) => {
      if (value <= 0) {
        callback(new Error('Must be greater than 0'));
      } else {
        callback();
      }
    };
    return {
      loading: true,
      loadingDialog: false,
      realEstateList: [],
      dialogCreateSelling: false,
      dialogCreateDonating: false,
      realForm: {
        price: 0,
        salePeriod: 0,
      },
      rules: {
        price: [
          { validator: checkArea, trigger: 'blur' },
        ],
        salePeriod: [
          { validator: checkArea, trigger: 'blur' },
        ],
      },
      DonatingForm: {
        proprietor: '',
      },
      rulesDonating: {
        proprietor: [
          { required: true, message: 'Select an owner', trigger: 'change' },
        ],
      },
      accountList: [],
      valItem: {},
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
    if (this.roles[0] === 'admin') {
      queryRealEstateList().then(response => {
        if (response !== null) {
          this.realEstateList = response;
        }
        this.loading = false;
      }).catch(_ => {
        this.loading = false;
      });
    } else {
      queryRealEstateList({ proprietor: this.accountId }).then(response => {
        if (response !== null) {
          this.realEstateList = response;
        }
        this.loading = false;
      }).catch(_ => {
        this.loading = false;
      });
    }
  },
  methods: {
    openDialog(item) {
      this.dialogCreateSelling = true;
      this.valItem = item;
    },
    openDonatingDialog(item) {
      this.dialogCreateDonating = true;
      this.valItem = item;
      queryAccountList().then(response => {
        if (response !== null) {
          // Filter out the admin and current user
          this.accountList = response.filter(item =>
            item.userName !== 'Admin' && item.accountId !== this.accountId
          );
        }
      });
    },
    createSelling(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('Create sale now?', 'Confirmation', {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            type: 'success',
          }).then(() => {
            this.loadingDialog = true;
            createSelling({
              objectOfSale: this.valItem.realEstateId,
              seller: this.valItem.proprietor,
              price: this.realForm.price,
              salePeriod: this.realForm.salePeriod,
            }).then(response => {
              this.loadingDialog = false;
              this.dialogCreateSelling = false;
              if (response !== null) {
                this.$message({
                  type: 'success',
                  message: 'Sale successful!',
                });
              } else {
                this.$message({
                  type: 'error',
                  message: 'Sale failed!',
                });
              }
              setTimeout(() => {
                window.location.reload();
              }, 1000);
            }).catch(_ => {
              this.loadingDialog = false;
              this.dialogCreateSelling = false;
            });
          }).catch(() => {
            this.loadingDialog = false;
            this.dialogCreateSelling = false;
            this.$message({
              type: 'info',
              message: 'Sale canceled',
            });
          });
        } else {
          return false;
        }
      });
    },
    createDonating(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('Create donation now?', 'Confirmation', {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            type: 'success',
          }).then(() => {
            this.loadingDialog = true;
            createDonating({
              objectOfDonating: this.valItem.realEstateId,
              donor: this.valItem.proprietor,
              grantee: this.DonatingForm.proprietor,
            }).then(response => {
              this.loadingDialog = false;
              this.dialogCreateDonating = false;
              if (response !== null) {
                this.$message({
                  type: 'success',
                  message: 'Donation successful!',
                });
              } else {
                this.$message({
                  type: 'error',
                  message: 'Donation failed!',
                });
              }
              setTimeout(() => {
                window.location.reload();
              }, 1000);
            }).catch(_ => {
              this.loadingDialog = false;
              this.dialogCreateDonating = false;
            });
          }).catch(() => {
            this.loadingDialog = false;
            this.dialogCreateDonating = false;
            this.$message({
              type: 'info',
              message: 'Donation canceled',
            });
          });
        } else {
          return false;
        }
      });
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    },
    selectGet(accountId) {
      this.DonatingForm.proprietor = accountId;
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

  .realEstate-card {
    width: 280px;
    height: 340px;
    margin: 18px;
  }
</style>
