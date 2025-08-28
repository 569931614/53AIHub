<template>
  <el-dialog
    class="service-dialog"
    :title="title || $t('service.title')"
    width="520px"
    :before-close="handleClose"
    :model-value="visible"
  >
    <div class="consult">
      <div class="consult-info">
        <div class="w-[70px] h-[70px] rounded-full bg-white flex justify-center items-center">
          <img class="w-[50px] h-[50px]" src="https://chat.53ai.com/images/extension_icon.png" />
        </div>
        <div class="consult-info--rt">
          <p class="name">
            {{ $t('service.consult_name') }}
          </p>
          <div class="flex-center">
            <ElIcon :size="16">
              <PhoneFilled />
            </ElIcon>
            <span class="mobile-number">186 8888 1185</span>
          </div>
        </div>
      </div>
      <div class="consult-ft px-[50px] py-[30px] max-md:px-4">
        <p
          class="max-w-64"
          v-html="
            `
              <span style='font-size: 15px;'>“</span>
              ${$t('service.desc_3') + $t('version.upgrade_renew')}
              <span style='font-size: 15px;'>”</span>
            `
          "
        ></p>
        <!-- <p v-html="`<span style='font-size: 15px;'>“</span>${$t('service.desc_1')}<br />${$t('service.desc_2')}<span style='font-size: 15px;'>”</span>`" /> -->
      </div>
      <div class="consult-img right-[30px] max-md:right-[15px]">
        <img src="https://hub.53ai.com/console/images/upgrade-qrcode.png" />
      </div>
    </div>
  </el-dialog>
</template>

<script lang="ts" setup>
import { PhoneFilled } from '@element-plus/icons-vue';

interface Emit {
  (e: 'update:visible', visible: boolean): void
}
interface Props {
  title?: string
  visible: boolean
}
withDefaults(defineProps<Props>(), {
  title: '',
  visible: false,
})

const emit = defineEmits<Emit>()

const handleClose = () => {
  emit('update:visible', false)
}
</script>

<style>
.service-dialog {
  border-radius: 10px 10px 11px 11px !important;
  background-color: #efefef !important;
  padding: 0 !important;
}

.service-dialog .el-dialog__header {
  padding: var(--el-dialog-padding-primary) !important;
}

.service-dialog .el-dialog__headerbtn {
  width: max-content !important;
  height: max-content !important;
}

.service-dialog .el-dialog__body {
  padding: 30px 0 0 !important;
  border: none !important;
  margin: 0 !important;
}
</style>

<style lang="scss" scoped>
.consult {
  &-info {
    display: flex;
    margin-bottom: 30px;
    padding: 0 38px;

    &--rt {
      margin-left: 20px;

      .name {
        font-size: 24px;
        font-weight: 500;
        vertical-align: text-bottom;
        margin-bottom: 10px;
      }

      .phone-icon {
        font-size: 16px;
        color: #0082f0;
      }

      .mobile-number {
        color: #666;
        margin-left: 5px;
        font-size: 18px;
      }
    }
  }

  &-ft {
    box-sizing: border-box;
    height: 125px;
    background-color: #0082f0;
    border-radius: 0 0 10px 10px;
    color: #fff;
    position: relative;

    &::before {
      content: '';
      width: 20px;
      height: 20px;
      background-color: #0082f0;
      transform: rotate(45deg);
      position: absolute;
      top: -10px;
      left: 125px;
    }

    p {
      line-height: 24px;
      position: relative;
      // &::before {
      //   content: '“';
      //   font-size: 40px;
      //   position: absolute;
      //   top: 10px;
      //   left: -20px;
      // }
      // &::after {
      //   content: '”';
      //   font-size: 40px;
      //   position: absolute;
      //   bottom: -15px;
      // }
    }
  }

  &-img {
    width: 104px;
    height: 104px;
    background: #fff;
    padding: 2px;
    box-shadow: 1px 1px 2px 0 rgb(7 116 208 / 20%);
    border-radius: 4px;
    position: absolute;
    bottom: 50px;

    img {
      width: 100%;
      height: 100%;
    }
  }
}
</style>
