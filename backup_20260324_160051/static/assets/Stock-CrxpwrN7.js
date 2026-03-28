const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/index-B4OaPzYO.js","assets/vue-vendor-CyCcT3D-.js","assets/element-plus-3gbl9L3C.js","assets/vendor-D_r0Q8NM.js","assets/pdf-export-CD-xqFA0.js","assets/index-BQaY74YA.css"])))=>i.map(i=>d[i]);
import{_ as fe}from"./pdf-export-CD-xqFA0.js";import{K as F,o as w,c as D,h as o,u as e,_ as r,$ as s,g as C,W as v,e as W,F as de,A as O,r as $,y as x,z as ie,N as Z,L as Ae,b as Ue,as as Ie,i as Oe,d as me}from"./vue-vendor-CyCcT3D-.js";import{j as se,r as re,k as J,i as Me,u as Re}from"./index-B4OaPzYO.js";import{a4 as ve,E as T,a5 as Ee,a6 as _e,T as Fe,M as We,R as He,k as Be,O as Ye,v as Ke}from"./element-plus-3gbl9L3C.js";import{D as ee}from"./Dialog-CQTxEK6t.js";import{T as Ge}from"./TableToolbar-P5duQHZY.js";import{W as ge,a as be}from"./WorkflowHistory-s-GMVy6O.js";import{_ as te}from"./_plugin-vue_export-helper-DlAUqK2U.js";import{P as Je}from"./ProjectSelector-D1gsdQC2.js";import"./vendor-D_r0Q8NM.js";const Qe={style:{"margin-top":"20px"}},Xe={__name:"InboundDetailDialog",props:{modelValue:{type:Boolean,default:!1},orderNo:{type:String,default:""}},emits:["update:modelValue"],setup(G,{emit:H}){const y=G,L=H,_=x(!1),f=x(!1),t=x({order_no:"",supplier:"",contact:"",receiver:"",inbound_date:"",remark:"",items:[],status:"",project_id:null,project_name:"",creator_name:"",created_at:"",updated_at:""}),N=x([]),c=i=>Number(i).toLocaleString("zh-CN",{minimumFractionDigits:2,maximumFractionDigits:2}),A=i=>({pending:"info",approved:"warning",rejected:"danger",completed:"success"})[i]||"info",h=i=>({pending:"待审核",approved:"已批准",rejected:"已拒绝",completed:"已完成"})[i]||i,B=i=>({pending:"入库单等待审核",approved:"入库单已批准，等待入库",rejected:"入库单已被拒绝",completed:"入库单已完成入库"})[i]||"",M=async i=>{try{const u=await se.getWorkflowHistory(i);N.value=u.data||[]}catch(u){console.error("获取工作流历史失败:",u)}},U=async()=>{if(y.orderNo){f.value=!0;try{const u=(await se.getList({pageSize:1e3})).data.find(k=>k.order_no===y.orderNo||k.inbound_no===y.orderNo);if(!u){T.error("未找到该入库单"),R();return}const m=await se.getDetail(u.id);t.value=m.data||m,await M(u.id)}catch(i){console.error("加载入库单详情失败:",i),T.error("加载入库单详情失败"),R()}finally{f.value=!1}}},R=()=>{_.value=!1,L("update:modelValue",!1)},g=()=>{const i=`
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <title>入库单 - ${t.value.order_no}</title>
      <style>
        * {
          margin: 0;
          padding: 0;
          box-sizing: border-box;
        }

        body {
          font-family: "Microsoft YaHei", Arial, sans-serif;
          font-size: 12px;
          line-height: 1.5;
          color: #333;
          padding: 20px;
        }

        .header {
          text-align: center;
          margin-bottom: 20px;
          border-bottom: 2px solid #333;
          padding-bottom: 10px;
        }

        .header h1 {
          font-size: 20px;
          margin-bottom: 5px;
        }

        .info-section {
          margin-bottom: 20px;
        }

        .info-row {
          display: flex;
          margin-bottom: 8px;
          border-bottom: 1px solid #eee;
          padding-bottom: 5px;
        }

        .info-label {
          font-weight: bold;
          width: 100px;
          flex-shrink: 0;
        }

        .info-value {
          flex: 1;
        }

        .section-title {
          font-size: 14px;
          font-weight: bold;
          margin: 20px 0 10px 0;
          padding-bottom: 5px;
          border-bottom: 1px solid #333;
        }

        table {
          width: 100%;
          border-collapse: collapse;
          margin-bottom: 20px;
        }

        table th,
        table td {
          border: 1px solid #333;
          padding: 8px;
          text-align: left;
        }

        table th {
          background-color: #f5f5f5;
          font-weight: bold;
          text-align: center;
        }

        table td {
          text-align: center;
        }

        .text-left {
          text-align: left !important;
        }

        .text-right {
          text-align: right !important;
        }

        @page {
          size: A4;
          margin: 15mm;
        }

        /* 防止表格行跨页 */
        tr {
          page-break-inside: avoid;
        }

        /* 防止元素内部跨页 */
        .info-section,
        .section-title {
          page-break-inside: avoid;
        }

        /* 避免在表格后立即分页 */
        table {
          page-break-inside: avoid;
        }
      </style>
    </head>
    <body>
      <div class="header">
        <h1>入库单</h1>
      </div>

      <div class="info-section">
        <div class="info-row">
          <span class="info-label">入库单号：</span>
          <span class="info-value">${t.value.order_no||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">状态：</span>
          <span class="info-value">${h(t.value.status)}</span>
        </div>
        <div class="info-row">
          <span class="info-label">供应商：</span>
          <span class="info-value">${t.value.supplier||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">联系人：</span>
          <span class="info-value">${t.value.contact||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">验收人：</span>
          <span class="info-value">${t.value.receiver||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">创建人：</span>
          <span class="info-value">${t.value.creator_name||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">创建时间：</span>
          <span class="info-value">${t.value.created_at||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">关联项目：</span>
          <span class="info-value">${t.value.project_name||"-"}</span>
        </div>
        ${t.value.remark?`
        <div class="info-row">
          <span class="info-label">备注：</span>
          <span class="info-value">${t.value.remark}</span>
        </div>
        `:""}
      </div>

      <div class="section-title">物资明细</div>
      <table>
        <thead>
          <tr>
            <th>物资名称</th>
            <th>规格型号</th>
            <th>材质</th>
            <th>单位</th>
            <th>数量</th>
            <th>单价</th>
            <th>金额</th>
          </tr>
        </thead>
        <tbody>
          ${t.value.items.map(m=>`
            <tr>
              <td class="text-left">${m.material_name||"-"}</td>
              <td>${m.spec||"-"}</td>
              <td>${m.material||"-"}</td>
              <td>${m.unit||"-"}</td>
              <td>${m.quantity||0}</td>
              <td>${m.unit_price?Number(m.unit_price).toFixed(2):"-"}</td>
              <td class="text-right">${((m.quantity||0)*(m.unit_price||0)).toFixed(2)}</td>
            </tr>
          `).join("")}
        </tbody>
      </table>

      <div style="margin-top: 30px; text-align: right; font-size: 10px; color: #999;">
        打印时间：${new Date().toLocaleString("zh-CN")}
      </div>
    </body>
    </html>
  `,u=window.open("","_blank");u?(u.document.write(i),u.document.close(),u.onload=()=>{u.print(),u.close()}):T.error("无法打开打印窗口，请检查浏览器设置")};return F(()=>y.modelValue,i=>{_.value=i,i&&U()}),F(_,i=>{i||L("update:modelValue",!1)}),(i,u)=>{const m=$("el-button"),k=$("el-tag"),S=$("el-descriptions-item"),d=$("el-descriptions"),p=$("el-divider"),q=$("el-table-column"),E=$("el-table");return w(),D(ee,{modelValue:_.value,"onUpdate:modelValue":u[0]||(u[0]=j=>_.value=j),title:"入库单详情",width:"900px",loading:f.value,onCancel:R},{extra:o(()=>[e(m,{type:"primary",icon:O(ve),onClick:g,size:"small"},{default:o(()=>[...u[1]||(u[1]=[r(" 打印 ",-1)])]),_:1},8,["icon"])]),default:o(()=>[e(d,{column:2,border:""},{default:o(()=>[e(S,{label:"入库单号",span:2},{default:o(()=>[e(k,{type:"primary"},{default:o(()=>[r(s(t.value.order_no||"-"),1)]),_:1})]),_:1}),e(S,{label:"状态"},{default:o(()=>[e(k,{type:A(t.value.status),size:"small"},{default:o(()=>[r(s(h(t.value.status)),1)]),_:1},8,["type"])]),_:1}),e(S,{label:"供应商"},{default:o(()=>[r(s(t.value.supplier||"-"),1)]),_:1}),e(S,{label:"联系人"},{default:o(()=>[r(s(t.value.contact||"-"),1)]),_:1}),e(S,{label:"验收人"},{default:o(()=>[r(s(t.value.receiver||"-"),1)]),_:1}),e(S,{label:"创建人"},{default:o(()=>[r(s(t.value.creator_name||"-"),1)]),_:1}),e(S,{label:"创建时间"},{default:o(()=>[r(s(t.value.created_at||"-"),1)]),_:1}),e(S,{label:"关联项目",span:2},{default:o(()=>[r(s(t.value.project_name||"-"),1)]),_:1}),t.value.remark?(w(),D(S,{key:0,label:"备注",span:2},{default:o(()=>[r(s(t.value.remark||"-"),1)]),_:1})):C("",!0)]),_:1}),v("div",Qe,[t.value.status?(w(),D(ge,{key:0,status:t.value.status,"status-time":t.value.updated_at||t.value.created_at,"status-description":B(t.value.status),"workflow-type":"inbound"},null,8,["status","status-time","status-description"])):C("",!0)]),e(p,{"content-position":"left"},{default:o(()=>[...u[2]||(u[2]=[r("物资明细",-1)])]),_:1}),e(E,{data:t.value.items,border:"",stripe:"",style:{width:"100%"},size:"small"},{default:o(()=>[e(q,{prop:"material_name",label:"物资名称","min-width":"150","show-overflow-tooltip":""}),e(q,{prop:"spec",label:"规格型号",width:"120","show-overflow-tooltip":""}),e(q,{prop:"material",label:"材质",width:"100","show-overflow-tooltip":""}),e(q,{prop:"unit",label:"单位",width:"80"}),e(q,{prop:"quantity",label:"数量",width:"100",align:"right"}),e(q,{prop:"unit_price",label:"单价",width:"100",align:"right"},{default:o(j=>[r(s(j.row.unit_price?c(j.row.unit_price):"-"),1)]),_:1}),e(q,{label:"金额",width:"120",align:"right"},{default:o(j=>[r(s(c((j.row.quantity||0)*(j.row.unit_price||0))),1)]),_:1}),e(q,{prop:"remark",label:"备注","min-width":"120","show-overflow-tooltip":""})]),_:1},8,["data"]),N.value.length>0?(w(),W(de,{key:0},[e(p,{"content-position":"left"},{default:o(()=>[...u[3]||(u[3]=[r("审批历史",-1)])]),_:1}),e(be,{histories:N.value},null,8,["histories"])],64)):C("",!0)]),_:1},8,["modelValue","loading"])}}},Ze=te(Xe,[["__scopeId","data-v-0fd3f22d"]]),et={key:1},tt={style:{"margin-top":"20px"}},at={__name:"RequisitionDetailDialog",props:{modelValue:{type:Boolean,default:!1},requisitionNo:{type:String,default:""}},emits:["update:modelValue"],setup(G,{emit:H}){const y=G,L=H,_=x(!1),f=x(!1),t=x({requisition_no:"",applicant:"",applicant_name:"",department:"",project_id:null,project_name:"",requisition_date:"",purpose:"",urgent:!1,remark:"",items:[],items_count:0,status:"",approved_by:"",approved_at:"",issued_by:"",issued_at:"",created_at:"",updated_at:""}),N=x([]),c=g=>({pending:"info",approved:"warning",rejected:"danger",issued:"success"})[g]||"info",A=g=>({pending:"待审核",approved:"已批准",rejected:"已拒绝",issued:"已发放"})[g]||g,h=g=>({pending:"出库单等待审核",approved:"出库单已批准，等待发放",rejected:"出库单已被拒绝",issued:"出库单已发放完成"})[g]||"",B=async g=>{try{const i=await re.getWorkflowHistory(g);N.value=i.data||[]}catch(i){console.error("获取工作流历史失败:",i)}},M=async()=>{if(y.requisitionNo){f.value=!0;try{const i=(await re.getList({pageSize:1e3})).data.find(m=>m.requisition_no===y.requisitionNo);if(!i){T.error("未找到该出库单"),U();return}const u=await re.getDetail(i.id);t.value=u.data||u,await B(i.id)}catch(g){console.error("加载出库单详情失败:",g),T.error("加载出库单详情失败"),U()}finally{f.value=!1}}},U=()=>{_.value=!1,L("update:modelValue",!1)},R=()=>{var u;const g=`
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <title>出库单 - ${t.value.requisition_no}</title>
      <style>
        * {
          margin: 0;
          padding: 0;
          box-sizing: border-box;
        }

        body {
          font-family: "Microsoft YaHei", Arial, sans-serif;
          font-size: 12px;
          line-height: 1.5;
          color: #333;
          padding: 20px;
        }

        .header {
          text-align: center;
          margin-bottom: 20px;
          border-bottom: 2px solid #333;
          padding-bottom: 10px;
        }

        .header h1 {
          font-size: 20px;
          margin-bottom: 5px;
        }

        .info-section {
          margin-bottom: 20px;
        }

        .info-row {
          display: flex;
          margin-bottom: 8px;
          border-bottom: 1px solid #eee;
          padding-bottom: 5px;
        }

        .info-label {
          font-weight: bold;
          width: 100px;
          flex-shrink: 0;
        }

        .info-value {
          flex: 1;
        }

        .section-title {
          font-size: 14px;
          font-weight: bold;
          margin: 20px 0 10px 0;
          padding-bottom: 5px;
          border-bottom: 1px solid #333;
        }

        table {
          width: 100%;
          border-collapse: collapse;
          margin-bottom: 20px;
        }

        table th,
        table td {
          border: 1px solid #333;
          padding: 8px;
          text-align: left;
        }

        table th {
          background-color: #f5f5f5;
          font-weight: bold;
          text-align: center;
        }

        table td {
          text-align: center;
        }

        .text-left {
          text-align: left !important;
        }

        @page {
          size: A4;
          margin: 15mm;
        }

        /* 防止表格行跨页 */
        tr {
          page-break-inside: avoid;
        }

        /* 防止元素内部跨页 */
        .info-section,
        .section-title {
          page-break-inside: avoid;
        }

        /* 避免在表格后立即分页 */
        table {
          page-break-inside: avoid;
        }
      </style>
    </head>
    <body>
      <div class="header">
        <h1>出库单</h1>
      </div>

      <div class="info-section">
        <div class="info-row">
          <span class="info-label">出库单号：</span>
          <span class="info-value">${t.value.requisition_no||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">状态：</span>
          <span class="info-value">${A(t.value.status)}</span>
        </div>
        <div class="info-row">
          <span class="info-label">项目名称：</span>
          <span class="info-value">${t.value.project_name||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">申请人：</span>
          <span class="info-value">${t.value.applicant_name||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">部门：</span>
          <span class="info-value">${t.value.department||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">申请日期：</span>
          <span class="info-value">${t.value.requisition_date||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">用途：</span>
          <span class="info-value">${t.value.purpose||"-"}</span>
        </div>
        ${t.value.urgent?`
        <div class="info-row">
          <span class="info-label">紧急：</span>
          <span class="info-value">是</span>
        </div>
        `:""}
        ${t.value.approved_by?`
        <div class="info-row">
          <span class="info-label">审核人：</span>
          <span class="info-value">${t.value.approved_by}</span>
        </div>
        `:""}
        ${t.value.approved_at?`
        <div class="info-row">
          <span class="info-label">审核时间：</span>
          <span class="info-value">${t.value.approved_at}</span>
        </div>
        `:""}
        ${t.value.issued_by?`
        <div class="info-row">
          <span class="info-label">发货人：</span>
          <span class="info-value">${t.value.issued_by}</span>
        </div>
        `:""}
        ${t.value.issued_at?`
        <div class="info-row">
          <span class="info-label">发货时间：</span>
          <span class="info-value">${t.value.issued_at}</span>
        </div>
        `:""}
        ${t.value.remark?`
        <div class="info-row">
          <span class="info-label">备注：</span>
          <span class="info-value">${t.value.remark}</span>
        </div>
        `:""}
      </div>

      <div class="section-title">物资明细 (${t.value.items_count||((u=t.value.items)==null?void 0:u.length)||0})</div>
      <table>
        <thead>
          <tr>
            <th>材质</th>
            <th>物资名称</th>
            <th>规格型号</th>
            <th>单位</th>
            <th>申请数量</th>
            <th>批准数量</th>
          </tr>
        </thead>
        <tbody>
          ${t.value.items.map(m=>`
            <tr>
              <td>${m.material||"-"}</td>
              <td class="text-left">${m.material_name||"-"}</td>
              <td>${m.specification||"-"}</td>
              <td>${m.unit||"-"}</td>
              <td>${m.requested_quantity||0}</td>
              <td>${m.approved_quantity||0}</td>
            </tr>
          `).join("")}
        </tbody>
      </table>

      <div style="margin-top: 30px; text-align: right; font-size: 10px; color: #999;">
        打印时间：${new Date().toLocaleString("zh-CN")}
      </div>
    </body>
    </html>
  `,i=window.open("","_blank");i?(i.document.write(g),i.document.close(),i.onload=()=>{i.print(),i.close()}):T.error("无法打开打印窗口，请检查浏览器设置")};return F(()=>y.modelValue,g=>{_.value=g,g&&M()}),F(_,g=>{g||L("update:modelValue",!1)}),(g,i)=>{const u=$("el-button"),m=$("el-tag"),k=$("el-descriptions-item"),S=$("el-descriptions"),d=$("el-divider"),p=$("el-table-column"),q=$("el-table");return w(),D(ee,{modelValue:_.value,"onUpdate:modelValue":i[0]||(i[0]=E=>_.value=E),title:"出库单详情",width:"900px",loading:f.value,onCancel:U},{extra:o(()=>[e(u,{type:"primary",icon:O(ve),onClick:R,size:"small"},{default:o(()=>[...i[1]||(i[1]=[r(" 打印 ",-1)])]),_:1},8,["icon"])]),default:o(()=>[e(S,{column:2,border:""},{default:o(()=>[e(k,{label:"出库单号",span:2},{default:o(()=>[e(m,{type:"primary"},{default:o(()=>[r(s(t.value.requisition_no||"-"),1)]),_:1})]),_:1}),e(k,{label:"状态"},{default:o(()=>[e(m,{type:c(t.value.status),size:"small"},{default:o(()=>[r(s(A(t.value.status)),1)]),_:1},8,["type"])]),_:1}),e(k,{label:"紧急"},{default:o(()=>[t.value.urgent?(w(),D(m,{key:0,type:"danger",size:"small"},{default:o(()=>[...i[2]||(i[2]=[r("紧急",-1)])]),_:1})):(w(),W("span",et,"否"))]),_:1}),e(k,{label:"项目名称",span:2},{default:o(()=>[r(s(t.value.project_name||"-"),1)]),_:1}),e(k,{label:"申请人"},{default:o(()=>[r(s(t.value.applicant_name||"-"),1)]),_:1}),e(k,{label:"部门"},{default:o(()=>[r(s(t.value.department||"-"),1)]),_:1}),e(k,{label:"申请日期"},{default:o(()=>[r(s(t.value.requisition_date||"-"),1)]),_:1}),e(k,{label:"用途"},{default:o(()=>[r(s(t.value.purpose||"-"),1)]),_:1}),e(k,{label:"创建时间"},{default:o(()=>[r(s(t.value.created_at||"-"),1)]),_:1}),t.value.approved_by?(w(),D(k,{key:0,label:"审核人"},{default:o(()=>[r(s(t.value.approved_by||"-"),1)]),_:1})):C("",!0),t.value.approved_at?(w(),D(k,{key:1,label:"审核时间"},{default:o(()=>[r(s(t.value.approved_at||"-"),1)]),_:1})):C("",!0),t.value.issued_by?(w(),D(k,{key:2,label:"发货人"},{default:o(()=>[r(s(t.value.issued_by||"-"),1)]),_:1})):C("",!0),t.value.issued_at?(w(),D(k,{key:3,label:"发货时间"},{default:o(()=>[r(s(t.value.issued_at||"-"),1)]),_:1})):C("",!0),e(k,{label:"备注",span:2},{default:o(()=>[r(s(t.value.remark||"-"),1)]),_:1})]),_:1}),v("div",tt,[t.value.status?(w(),D(ge,{key:0,status:t.value.status,"status-time":t.value.updated_at||t.value.created_at,"status-description":h(t.value.status),"workflow-type":"requisition"},null,8,["status","status-time","status-description"])):C("",!0)]),e(d,{"content-position":"left"},{default:o(()=>{var E;return[r("物资明细 ("+s(t.value.items_count||((E=t.value.items)==null?void 0:E.length)||0)+")",1)]}),_:1}),e(q,{data:t.value.items,border:"",stripe:"",style:{width:"100%"},size:"small"},{default:o(()=>[e(p,{prop:"material",label:"材质",width:"100","show-overflow-tooltip":""}),e(p,{prop:"material_name",label:"物资名称","min-width":"150","show-overflow-tooltip":""}),e(p,{prop:"specification",label:"规格型号","min-width":"150","show-overflow-tooltip":""}),e(p,{prop:"unit",label:"单位",width:"80"}),e(p,{prop:"requested_quantity",label:"申请数量",width:"100",align:"right"}),e(p,{prop:"approved_quantity",label:"批准数量",width:"100",align:"right"})]),_:1},8,["data"]),N.value.length>0?(w(),W(de,{key:0},[e(d,{"content-position":"left"},{default:o(()=>[...i[3]||(i[3]=[r("审批历史",-1)])]),_:1}),e(be,{histories:N.value},null,8,["histories"])],64)):C("",!0)]),_:1},8,["modelValue","loading"])}}},ot=te(at,[["__scopeId","data-v-5e654ede"]]),lt={class:"stock-dialog-content"},nt={class:"material-section"},it={key:0,class:"material-card"},st={class:"material-card-header"},rt={class:"material-title"},dt={class:"material-name"},ut={class:"material-card-body"},pt={class:"info-grid"},ct={class:"info-item"},ft={class:"info-value"},mt={class:"info-item"},vt={class:"info-value"},_t={class:"info-item"},gt={class:"info-value"},bt={key:0,class:"info-item"},ht={class:"info-value price"},yt={class:"stock-highlight"},wt={class:"stock-value"},kt={class:"operation-form"},xt={class:"form-row"},$t={class:"form-item-wrapper"},qt={class:"form-label"},Vt={key:0,class:"form-hint"},zt={key:0,class:"amount-bar"},Dt={class:"amount-info"},St={class:"amount-value"},Ct={class:"amount-detail"},jt={class:"detail-text"},Tt={__name:"StockOperationDialog",props:{modelValue:{type:Boolean,default:!1},materialId:{type:Number,default:null},operationType:{type:String,default:"in"},stockData:{type:Object,default:null}},emits:["update:modelValue","success"],setup(G,{emit:H}){const y=G,L=H,_=x(!1),f=x(!1),t=[{value:"in",label:"入库",icon:Ee,color:"#67c23a"},{value:"out",label:"出库",icon:_e,color:"#e6a23c"},{value:"adjust",label:"调整",icon:Fe,color:"#409eff"}],N=ie(()=>{const d=t.find(p=>p.value===c.type);return d?`${d.label}操作`:"库存操作"}),c=Z({material_id:null,quantity:null,price:null,type:"in",remark:""}),A=x([]),h=ie(()=>{const d=A.value.find(p=>p.id===c.material_id);return y.stockData?{id:y.stockData.material_id,code:y.stockData.material_code,name:y.stockData.material_name,category:y.stockData.category,specification:y.stockData.specification,unit:y.stockData.unit,quantity:y.stockData.quantity,safety_stock:y.stockData.safety_stock,price:d==null?void 0:d.price,...d}:d}),B=ie(()=>{var d;return c.quantity&&(c.price||((d=h.value)==null?void 0:d.price))}),M=d=>Number(d||0).toFixed(2),U=()=>{var q;const d=c.price||((q=h.value)==null?void 0:q.price)||0,p=c.quantity||0;return d*p},R=async()=>{try{const{data:d}=await Me.getList({pageSize:1e3});A.value=d||[]}catch(d){console.error("获取物资列表失败:",d)}},g=()=>{const d=h.value;d&&(c.price=d.price)},i=()=>{Object.assign(c,{material_id:null,quantity:null,price:null,type:y.operationType,remark:""})},u=()=>c.material_id?!c.quantity||c.quantity<=0?(T.warning("请输入有效的数量"),!1):c.type==="out"&&h.value&&c.quantity>h.value.quantity?(T.warning(`出库数量不能超过当前库存 ${h.value.quantity}`),!1):!0:(T.warning("请选择物资"),!1),m=async()=>{var p;if(!u())return;const d=(p=y.stockData)==null?void 0:p.id;if(!d){T.error("无法获取库存记录ID");return}try{f.value=!0;const q={quantity:c.quantity,remark:c.remark};c.type==="in"||c.type==="adjust"?(await J.in(d,q),T.success("入库成功")):(await J.out(d,q),T.success("出库成功")),L("success"),k()}catch(q){console.error("提交失败:",q),T.error(q.message||"操作失败")}finally{f.value=!1}},k=()=>{_.value=!1,L("update:modelValue",!1)},S=async()=>{i(),await R(),y.materialId&&(c.material_id=y.materialId,c.type=y.operationType,g())};return F(()=>y.modelValue,d=>{_.value=d,d&&S()}),F(_,d=>{d||L("update:modelValue",!1)}),(d,p)=>{const q=$("el-tag"),E=$("el-input-number");return w(),D(ee,{modelValue:_.value,"onUpdate:modelValue":p[1]||(p[1]=j=>_.value=j),title:N.value,width:"580px",loading:f.value,onConfirm:m,onCancel:k},{default:o(()=>{var j,Q;return[v("div",lt,[v("div",nt,[h.value?(w(),W("div",it,[v("div",st,[v("div",rt,[v("span",dt,s(h.value.name),1),e(q,{size:"small",type:"info"},{default:o(()=>[r(s(h.value.category||"-"),1)]),_:1})])]),v("div",ut,[v("div",pt,[v("div",ct,[p[2]||(p[2]=v("span",{class:"info-label"},"编码",-1)),v("span",ft,s(h.value.code||"-"),1)]),v("div",mt,[p[3]||(p[3]=v("span",{class:"info-label"},"规格",-1)),v("span",vt,s(h.value.specification||"-"),1)]),v("div",_t,[p[4]||(p[4]=v("span",{class:"info-label"},"单位",-1)),v("span",gt,s(h.value.unit||"-"),1)]),h.value.price?(w(),W("div",bt,[p[5]||(p[5]=v("span",{class:"info-label"},"单价",-1)),v("span",ht,"¥"+s(M(h.value.price)),1)])):C("",!0)]),v("div",yt,[p[6]||(p[6]=v("span",{class:"stock-label"},"当前库存",-1)),v("span",wt,s(h.value.quantity||0)+" "+s(h.value.unit),1)])])])):C("",!0)]),v("div",kt,[v("div",xt,[v("div",$t,[v("label",qt,[p[7]||(p[7]=v("span",null,"数量",-1)),c.type==="out"&&h.value?(w(),W("span",Vt," 最大可出库: "+s(h.value.quantity),1)):C("",!0)]),e(E,{modelValue:c.quantity,"onUpdate:modelValue":p[0]||(p[0]=ae=>c.quantity=ae),min:1,max:c.type==="out"?(j=h.value)==null?void 0:j.quantity:void 0,step:1,precision:0,disabled:!c.material_id,class:"full-width-input",placeholder:"请输入数量"},null,8,["modelValue","max","disabled"])])]),B.value?(w(),W("div",zt,[v("div",Dt,[p[8]||(p[8]=v("span",{class:"amount-label"},"预计金额",-1)),v("span",St,"¥"+s(M(U())),1)]),v("div",Ct,[v("span",jt,s(c.quantity)+" × ¥"+s(c.price||((Q=h.value)==null?void 0:Q.price)||0),1)])])):C("",!0)])])]}),_:1},8,["modelValue","title","loading"])}}},Lt=te(Tt,[["__scopeId","data-v-23522ac3"]]),Nt={class:"stock-container"},Pt={key:1},At={__name:"Stock",setup(G){const H=Re(),y=x(!1),L=x([]),_=Z({page:1,pageSize:20,total:0}),f=Z({keyword:"",category:"",project_id:"",status:""}),t=x([]),N=x([]),c=x(!1),A=x(null),h=x("in"),B=x(null),M=x(!1),U=x(!1),R=x([]),g=Z({page:1,pageSize:20,total:0}),i=x(null),u=x(!1),m=x(!1),k=x(""),S=x(""),d=async()=>{y.value=!0;try{let a=[];f.project_id&&(a=p(f.project_id,t.value));const l={page:_.page,page_size:_.pageSize,search:f.keyword||void 0,category:f.category||void 0,project_ids:a.length>0?a.join(","):void 0,status:f.status||void 0},{data:b,pagination:V}=await J.getList(l);L.value=b||[],_.total=(V==null?void 0:V.total)||0}catch(a){console.error("获取库存列表失败:",a)}finally{y.value=!1}},p=(a,l)=>{const b=[a],V=P=>{for(const I of P){if(I.id===a){const Y=z=>{if(z.children&&z.children.length>0)for(const K of z.children)b.push(K.id),Y(K)};return Y(I),!0}if(I.children&&I.children.length>0&&V(I.children))return!0}return!1};return V(l),b},q=async()=>{try{const{projectApi:a}=await fe(async()=>{const{projectApi:b}=await import("./index-B4OaPzYO.js").then(V=>V.t);return{projectApi:b}},__vite__mapDeps([0,1,2,3,4,5])),{data:l}=await a.getList({pageSize:1e3});t.value=E(l||[])}catch(a){console.error("获取项目列表失败:",a)}},E=a=>{if(!a||a.length===0)return[];const l=new Map;a.forEach(V=>{l.set(V.id,{...V,children:[]})});const b=[];return a.forEach(V=>{const P=l.get(V.id);if(!V.parent_id)b.push(P);else{const I=l.get(V.parent_id);I?I.children.push(P):b.push(P)}}),b},j=async()=>{try{const{materialApi:a}=await fe(async()=>{const{materialApi:b}=await import("./index-B4OaPzYO.js").then(V=>V.t);return{materialApi:b}},__vite__mapDeps([0,1,2,3,4,5])),{data:l}=await a.getCategories();N.value=l||[]}catch(a){console.error("获取物资分类列表失败:",a)}},Q=()=>{_.page=1,d()},ae=()=>{f.keyword="",f.category="",f.project_id="",f.status="",_.page=1,d()},he=a=>{A.value=a.material_id,h.value="in",B.value=a,c.value=!0},ye=a=>{A.value=a.material_id,h.value="out",B.value=a,c.value=!0},we=a=>{i.value=a.id,M.value=!0,oe()},oe=async()=>{if(i.value){U.value=!0;try{const a={page:g.page,page_size:g.pageSize,stock_id:i.value},{data:l,pagination:b}=await J.getLogs(a);R.value=l||[],g.total=(b==null?void 0:b.total)||0}catch(a){console.error("获取库存日志失败:",a)}finally{U.value=!1}}},ke=async()=>{try{const a=await J.export(f),l=new Blob([a],{type:"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}),b=window.URL.createObjectURL(l),V=document.createElement("a");V.href=b,V.download=`库存列表_${new Date().getTime()}.xlsx`,V.click(),window.URL.revokeObjectURL(b),T.success("导出成功")}catch(a){console.error("导出失败:",a)}},xe=a=>Number(a).toLocaleString("zh-CN",{minimumFractionDigits:2,maximumFractionDigits:2}),$e=a=>a<=0?"danger":a<10?"warning":"success",qe=a=>{const l=a.quantity||0,b=a.safe_stock||0;return l<=0?"danger":l<b?"warning":"success"},Ve=a=>{const l=a.quantity||0,b=a.safe_stock||0;return l<=0?"库存不足":l<b?"库存偏低":"正常"},ze=a=>({in:"success",out:"warning",adjust:"info"})[a]||"info",De=a=>({in:"入库",out:"出库",adjust:"调整"})[a]||a,Se=a=>a?a.includes("入库单")||a.includes("出库单"):!1,Ce=a=>{if(a.inbound_code)return`入库单-${a.inbound_code}`;if(a.requisition_code)return`出库单-${a.requisition_code}`;if(!a.remark)return"-";const l=a.remark.match(/入库单(?:审核入库)?\s*[:：]?\s*(\w+)/);if(l)return`入库单-${l[1]}`;const b=a.remark.match(/出库单(?:发放)?\s*[:：]?\s*(\w+)/);return b?`出库单-${b[1]}`:a.remark},je=a=>{if(a.inbound_code){k.value=a.inbound_code,u.value=!0;return}if(a.requisition_code){S.value=a.requisition_code,m.value=!0;return}if(!a.remark)return;const l=a.remark.match(/入库单(?:审核入库)?\s*[:：]?\s*(\w+)/);if(l){k.value=l[1],u.value=!0;return}const b=a.remark.match(/出库单(?:发放)?\s*[:：]?\s*(\w+)/);if(b){S.value=b[1],m.value=!0;return}};let le=null;const X=()=>{le&&clearTimeout(le),le=setTimeout(()=>{_.page=1,d()},500)},Te=a=>{_.page=a,d()},Le=a=>{_.pageSize=a,_.page=1,d()};return F(()=>f.keyword,X),F(()=>f.category,X),F(()=>f.project_id,X),F(()=>f.status,X),Ae(()=>{q(),j(),d()}),(a,l)=>{const b=$("el-icon"),V=$("el-input"),P=$("el-option"),I=$("el-select"),Y=$("el-button"),z=$("el-table-column"),K=$("el-tag"),ue=$("el-table"),pe=$("el-pagination"),Ne=$("el-card"),Pe=$("el-link"),ce=Ue("loading");return w(),W("div",Nt,[e(Ne,{shadow:"never"},{default:o(()=>[e(Ge,null,{left:o(()=>[e(Je,{modelValue:f.project_id,"onUpdate:modelValue":l[0]||(l[0]=n=>f.project_id=n),projects:t.value,placeholder:"选择项目（支持层级显示）",width:"300px"},null,8,["modelValue","projects"]),e(V,{modelValue:f.keyword,"onUpdate:modelValue":l[1]||(l[1]=n=>f.keyword=n),placeholder:"搜索物资名称、编码",clearable:"",style:{width:"250px"},onKeyup:Ie(Q,["enter"])},{prefix:o(()=>[e(b,null,{default:o(()=>[e(O(He))]),_:1})]),_:1},8,["modelValue"]),e(I,{modelValue:f.category,"onUpdate:modelValue":l[2]||(l[2]=n=>f.category=n),placeholder:"物资分类",clearable:"",style:{width:"150px"}},{default:o(()=>[e(P,{label:"全部",value:""}),(w(!0),W(de,null,Oe(N.value,n=>(w(),D(P,{key:n.id,label:n.name,value:n.name},null,8,["label","value"]))),128))]),_:1},8,["modelValue"]),e(I,{modelValue:f.status,"onUpdate:modelValue":l[3]||(l[3]=n=>f.status=n),placeholder:"库存状态",clearable:"",style:{width:"150px"}},{default:o(()=>[e(P,{label:"所有状态",value:""}),e(P,{label:"正常",value:"normal"}),e(P,{label:"库存偏低",value:"low"}),e(P,{label:"库存不足",value:"shortage"})]),_:1},8,["modelValue"]),e(Y,{icon:O(Be),onClick:ae},{default:o(()=>[...l[12]||(l[12]=[r("重置",-1)])]),_:1},8,["icon"])]),right:o(()=>[O(H).hasPermission("stock_export")?(w(),D(Y,{key:0,type:"success",icon:O(We),onClick:ke},{default:o(()=>[...l[13]||(l[13]=[r(" 导出 ",-1)])]),_:1},8,["icon"])):C("",!0)]),_:1}),me((w(),D(ue,{data:L.value,border:"",stripe:"",style:{width:"100%"}},{default:o(()=>[e(z,{prop:"material_code",label:"物资编码",width:"130"}),e(z,{prop:"material_name",label:"物资名称","min-width":"150","show-overflow-tooltip":""}),e(z,{prop:"category",label:"分类",width:"100"},{default:o(n=>[e(K,{size:"small"},{default:o(()=>[r(s(n.row.category||"-"),1)]),_:2},1024)]),_:1}),e(z,{prop:"specification",label:"规格型号",width:"120","show-overflow-tooltip":""}),e(z,{label:"材质",width:"100","show-overflow-tooltip":""},{default:o(n=>[r(s(n.row.material||"-"),1)]),_:1}),e(z,{prop:"unit",label:"单位",width:"80"}),e(z,{prop:"quantity",label:"库存数量",width:"120",align:"right"},{default:o(n=>[e(K,{type:$e(n.row.quantity),size:"large"},{default:o(()=>[r(s(n.row.quantity||0),1)]),_:2},1032,["type"])]),_:1}),e(z,{prop:"safety_stock",label:"安全库存",width:"100",align:"right"},{default:o(n=>[r(s(n.row.safety_stock||"-"),1)]),_:1}),e(z,{prop:"stock_status",label:"库存状态",width:"100"},{default:o(n=>[e(K,{type:qe(n.row),size:"small"},{default:o(()=>[r(s(Ve(n.row)),1)]),_:2},1032,["type"])]),_:1}),e(z,{prop:"project_name",label:"关联项目","min-width":"150","show-overflow-tooltip":""},{default:o(n=>[r(s(n.row.project_name||"-"),1)]),_:1}),e(z,{prop:"updated_at",label:"最后更新",width:"160"},{default:o(n=>[r(s(n.row.updated_at||"-"),1)]),_:1}),e(z,{label:"操作",width:"260",fixed:"right"},{default:o(n=>[O(H).hasPermission("stock_in")?(w(),D(Y,{key:0,type:"success",size:"small",icon:O(Ye),onClick:ne=>he(n.row)},{default:o(()=>[...l[14]||(l[14]=[r(" 入库 ",-1)])]),_:1},8,["icon","onClick"])):C("",!0),O(H).hasPermission("stock_out")?(w(),D(Y,{key:1,type:"warning",size:"small",icon:O(_e),onClick:ne=>ye(n.row)},{default:o(()=>[...l[15]||(l[15]=[r(" 出库 ",-1)])]),_:1},8,["icon","onClick"])):C("",!0),e(Y,{type:"primary",size:"small",icon:O(Ke),onClick:ne=>we(n.row)},{default:o(()=>[...l[16]||(l[16]=[r(" 日志 ",-1)])]),_:1},8,["icon","onClick"])]),_:1})]),_:1},8,["data"])),[[ce,y.value]]),e(pe,{"current-page":_.page,"onUpdate:currentPage":l[4]||(l[4]=n=>_.page=n),"page-size":_.pageSize,"onUpdate:pageSize":l[5]||(l[5]=n=>_.pageSize=n),"page-sizes":[10,20,50,100],total:_.total,layout:"total, sizes, prev, pager, next, jumper",onSizeChange:Le,onCurrentChange:Te,class:"mt-20"},null,8,["current-page","page-size","total"])]),_:1}),e(Lt,{modelValue:c.value,"onUpdate:modelValue":l[6]||(l[6]=n=>c.value=n),"material-id":A.value,"operation-type":h.value,"stock-data":B.value,onSuccess:d},null,8,["modelValue","material-id","operation-type","stock-data"]),e(ee,{modelValue:M.value,"onUpdate:modelValue":l[9]||(l[9]=n=>M.value=n),title:"库存日志",width:"900px","show-footer":!1},{default:o(()=>[me((w(),D(ue,{data:R.value,border:"",stripe:"","max-height":"400"},{default:o(()=>[e(z,{prop:"created_at",label:"时间",width:"160"}),e(z,{prop:"type",label:"类型",width:"80"},{default:o(n=>[e(K,{type:ze(n.row.type),size:"small"},{default:o(()=>[r(s(De(n.row.type)),1)]),_:2},1032,["type"])]),_:1}),e(z,{prop:"quantity",label:"数量",width:"80",align:"right"}),e(z,{prop:"quantity_before",label:"操作前",width:"80",align:"right"}),e(z,{prop:"quantity_after",label:"操作后",width:"80",align:"right"}),e(z,{prop:"price",label:"单价",width:"80",align:"right"},{default:o(n=>[r(s(n.row.price?xe(n.row.price):"-"),1)]),_:1}),e(z,{prop:"remark",label:"备注","min-width":"200"},{default:o(n=>[n.row.inbound_code||n.row.requisition_code||Se(n.row.remark)?(w(),D(Pe,{key:0,type:"primary",onClick:ne=>je(n.row)},{default:o(()=>[r(s(Ce(n.row)),1)]),_:2},1032,["onClick"])):(w(),W("span",Pt,s(n.row.remark||"-"),1))]),_:1})]),_:1},8,["data"])),[[ce,U.value]]),e(pe,{"current-page":g.page,"onUpdate:currentPage":l[7]||(l[7]=n=>g.page=n),"page-size":g.pageSize,"onUpdate:pageSize":l[8]||(l[8]=n=>g.pageSize=n),"page-sizes":[10,20,50],total:g.total,layout:"total, sizes, prev, pager, next",onSizeChange:oe,onCurrentChange:oe,class:"mt-20"},null,8,["current-page","page-size","total"])]),_:1},8,["modelValue"]),e(Ze,{modelValue:u.value,"onUpdate:modelValue":l[10]||(l[10]=n=>u.value=n),"order-no":k.value},null,8,["modelValue","order-no"]),e(ot,{modelValue:m.value,"onUpdate:modelValue":l[11]||(l[11]=n=>m.value=n),"requisition-no":S.value},null,8,["modelValue","requisition-no"])])}}},Yt=te(At,[["__scopeId","data-v-265ef92d"]]);export{Yt as default};
