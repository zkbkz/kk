<template>
    <div>
        <MonitorRouter />

        <LayoutContent v-loading="loading" :title="$t('menu.settings', 2)">
            <template #main>
                <el-form :model="form" @submit.prevent label-position="left" label-width="160px">
                    <el-row>
                        <el-col :span="1"><br /></el-col>
                        <el-col :span="12">
                            <h3 class="mb-4">{{ $t('monitor.hostMonitor') }}</h3>
                            <el-form-item :label="$t('monitor.enableMonitor')" prop="monitorStatus">
                                <el-switch
                                    v-permission
                                    v-node-admin
                                    @change="onSaveStatus('MonitorStatus', form.monitorStatus)"
                                    v-model="form.monitorStatus"
                                    active-value="Enable"
                                    inactive-value="Disable"
                                />
                            </el-form-item>
                            <el-form-item :label="$t('monitor.storeDays')" prop="monitorStoreDays">
                                <el-input disabled v-model="form.monitorStoreDays">
                                    <template #append>
                                        <el-button
                                            v-permission
                                            v-node-admin
                                            @click="onChangeStoreDays('MonitorStoreDays', form.monitorStoreDays)"
                                            icon="Setting"
                                        >
                                            {{ $t('commons.button.set') }}
                                        </el-button>
                                    </template>
                                </el-input>
                            </el-form-item>
                            <el-form-item :label="$t('monitor.interval')" prop="monitorIntervalItem">
                                <el-input disabled v-model="form.monitorIntervalItem">
                                    <template #append>
                                        <el-button
                                            v-permission
                                            v-node-admin
                                            @click="onChangeInterval('MonitorInterval', form.monitorInterval)"
                                            icon="Setting"
                                        >
                                            {{ $t('commons.button.set') }}
                                        </el-button>
                                    </template>
                                </el-input>
                            </el-form-item>
                            <el-form-item :label="$t('monitor.defaultNetwork')">
                                <el-input disabled v-model="form.defaultNetwork">
                                    <template #append>
                                        <el-button v-permission v-node-admin @click="onChangeNetwork" icon="Setting">
                                            {{ $t('commons.button.set') }}
                                        </el-button>
                                    </template>
                                </el-input>
                                <span class="input-help">{{ $t('monitor.defaultNetworkHelper') }}</span>
                            </el-form-item>
                            <el-form-item :label="$t('monitor.defaultIO')">
                                <el-input disabled v-model="form.defaultIO">
                                    <template #append>
                                        <el-button v-permission v-node-admin @click="onChangeIO" icon="Setting">
                                            {{ $t('commons.button.set') }}
                                        </el-button>
                                    </template>
                                </el-input>
                                <span class="input-help">{{ $t('monitor.defaultIOHelper') }}</span>
                            </el-form-item>
                            <el-form-item>
                                <el-button v-permission v-node-admin @click="onClean('host')" icon="Delete">
                                    {{ $t('monitor.cleanMonitor') }}
                                </el-button>
                            </el-form-item>
                        </el-col>
                    </el-row>
                </el-form>
            </template>
        </LayoutContent>

        <Interval ref="intervalRef" @search="search" />
        <StoreDays ref="daysRef" @search="search" />
        <Network ref="networkRef" @search="search()" />
        <IO ref="ioRef" @search="search()" />
    </div>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { ElMessageBox } from 'element-plus';
import { cleanMonitors, loadMonitorSetting, updateMonitorSetting } from '@/api/modules/host';
import MonitorRouter from '@/views/host/monitor/index.vue';
import Interval from '@/views/host/monitor/setting/interval/index.vue';
import StoreDays from '@/views/host/monitor/setting/days/index.vue';
import Network from '@/views/host/monitor/setting/default-network/index.vue';
import IO from '@/views/host/monitor/setting/default-io/index.vue';
import i18n from '@/lang';
import { MsgSuccess } from '@/utils/message';
import { splitTimeFromSecond, transTimeUnit } from '@/utils/validate';
const loading = ref();
const form = reactive({
    monitorStatus: 'Disable',
    monitorStoreDays: 30,
    monitorInterval: 300,
    monitorIntervalItem: '',
    defaultNetwork: '',
    defaultIO: '',
});

const intervalRef = ref();
const daysRef = ref();
const networkRef = ref();
const ioRef = ref();

const search = async () => {
    const res = await loadMonitorSetting();
    form.monitorStatus = res.data.monitorStatus;
    form.monitorInterval = Number(res.data.monitorInterval);
    const item = splitTimeFromSecond(form.monitorInterval);
    form.monitorIntervalItem = transTimeUnit(item.timeItem + item.timeUnit);

    form.monitorStoreDays = Number(res.data.monitorStoreDays);
    form.defaultNetwork =
        res.data.defaultNetwork === 'all' ? i18n.global.t('commons.table.all') : res.data.defaultNetwork;
    form.defaultIO = res.data.defaultIO === 'all' ? i18n.global.t('commons.table.all') : res.data.defaultIO;
};

const onSaveStatus = async (key: 'MonitorStatus', value: string) => {
    loading.value = true;
    await updateMonitorSetting(key, value)
        .then(() => {
            loading.value = false;
            MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
        })
        .catch(async () => {
            loading.value = false;
            await search();
        });
};

const onChangeStoreDays = (key: 'MonitorStoreDays', value: number) => {
    daysRef.value.acceptParams({ key, monitorStoreDays: value });
};
const onChangeInterval = (key: 'MonitorInterval', value: number) => {
    const interval = splitTimeFromSecond(value);
    intervalRef.value.acceptParams({ key, timeItem: interval.timeItem, timeUnit: interval.timeUnit });
};
const onChangeNetwork = () => {
    networkRef.value.acceptParams({ defaultNetwork: form.defaultNetwork });
};
const onChangeIO = () => {
    ioRef.value.acceptParams({ defaultIO: form.defaultIO });
};

const onClean = async () => {
    const name = i18n.global.t('monitor.hostMonitor');
    ElMessageBox.confirm(i18n.global.t('monitor.cleanHelper', [name]), i18n.global.t('monitor.cleanMonitor'), {
        confirmButtonText: i18n.global.t('commons.button.confirm'),
        cancelButtonText: i18n.global.t('commons.button.cancel'),
        type: 'info',
    }).then(async () => {
        loading.value = true;
        await cleanMonitors()
            .then(() => {
                loading.value = false;
                MsgSuccess(i18n.global.t('commons.msg.operationSuccess'));
            })
            .catch(() => {
                loading.value = false;
            });
    });
};

onMounted(() => {
    search();
});
</script>
