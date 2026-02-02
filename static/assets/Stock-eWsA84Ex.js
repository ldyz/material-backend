const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["assets/index-BAziDdcB.js","assets/vue-vendor-Ch5YuIb-.js","assets/element-plus-Dsy0BBrW.js","assets/index-Vjn9Jgty.css"])))=>i.map(i=>d[i]);
import{i as se,r as re,u as Qe,_ as _e,h as Q,m as Ge}from"./index-BAziDdcB.js";import{h as I,y as g,z as k,D as t,C as e,Q as s,R as d,P as q,E as G,B as Y,J as te,u as D,ab as p,r as b,j as Xe,c as ge,l as ee,ac as Ze,am as et,a5 as be,M as ye}from"./vue-vendor-Ch5YuIb-.js";import{ag as he,E as R,A as tt,ah as at,H as lt,I as ot,J as nt,j as it}from"./element-plus-Dsy0BBrW.js";import{D as ae}from"./Dialog-mEQRQXk5.js";import{T as st}from"./TableToolbar-DstB8QC1.js";import{W as we,a as xe}from"./WorkflowHistory-BBLhbZ5d.js";import{_ as de}from"./_plugin-vue_export-helper-DlAUqK2U.js";import{P as rt}from"./ProjectSelector-BNM1bkWN.js";const dt={style:{"margin-top":"20px"}},ut={__name:"InboundDetailDialog",props:{modelValue:{type:Boolean,default:!1},orderNo:{type:String,default:""}},emits:["update:modelValue"],setup(X,{emit:H}){const S=X,E=H,_=b(!1),c=b(!1),a=b({order_no:"",supplier:"",contact:"",receiver:"",inbound_date:"",remark:"",items:[],status:"",project_id:null,project_name:"",creator_name:"",created_at:"",updated_at:""}),T=b([]),L=i=>Number(i).toLocaleString("zh-CN",{minimumFractionDigits:2,maximumFractionDigits:2}),N=i=>({pending:"info",approved:"warning",rejected:"danger",completed:"success"})[i]||"info",J=i=>({pending:"待审核",approved:"已批准",rejected:"已拒绝",completed:"已完成"})[i]||i,F=i=>({pending:"入库单等待审核",approved:"入库单已批准，等待入库",rejected:"入库单已被拒绝",completed:"入库单已完成入库"})[i]||"",P=async i=>{try{const r=await se.getWorkflowHistory(i);T.value=r.data||[]}catch(r){console.error("获取工作流历史失败:",r)}},f=async()=>{if(S.orderNo){c.value=!0;try{const r=(await se.getList({pageSize:1e3})).data.find(h=>h.order_no===S.orderNo||h.inbound_no===S.orderNo);if(!r){R.error("未找到该入库单"),A();return}const u=await se.getDetail(r.id);a.value=u.data||u,await P(r.id)}catch(i){console.error("加载入库单详情失败:",i),R.error("加载入库单详情失败"),A()}finally{c.value=!1}}},A=()=>{_.value=!1,E("update:modelValue",!1)},v=()=>{const i=`
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <title>入库单 - ${a.value.order_no}</title>
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
          <span class="info-value">${a.value.order_no||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">状态：</span>
          <span class="info-value">${J(a.value.status)}</span>
        </div>
        <div class="info-row">
          <span class="info-label">供应商：</span>
          <span class="info-value">${a.value.supplier||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">联系人：</span>
          <span class="info-value">${a.value.contact||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">验收人：</span>
          <span class="info-value">${a.value.receiver||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">创建人：</span>
          <span class="info-value">${a.value.creator_name||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">创建时间：</span>
          <span class="info-value">${a.value.created_at||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">关联项目：</span>
          <span class="info-value">${a.value.project_name||"-"}</span>
        </div>
        ${a.value.remark?`
        <div class="info-row">
          <span class="info-label">备注：</span>
          <span class="info-value">${a.value.remark}</span>
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
          ${a.value.items.map(u=>`
            <tr>
              <td class="text-left">${u.material_name||"-"}</td>
              <td>${u.spec||"-"}</td>
              <td>${u.material||"-"}</td>
              <td>${u.unit||"-"}</td>
              <td>${u.quantity||0}</td>
              <td>${u.unit_price?Number(u.unit_price).toFixed(2):"-"}</td>
              <td class="text-right">${((u.quantity||0)*(u.unit_price||0)).toFixed(2)}</td>
            </tr>
          `).join("")}
        </tbody>
      </table>

      <div style="margin-top: 30px; text-align: right; font-size: 10px; color: #999;">
        打印时间：${new Date().toLocaleString("zh-CN")}
      </div>
    </body>
    </html>
  `,r=window.open("","_blank");r?(r.document.write(i),r.document.close(),r.onload=()=>{r.print(),r.close()}):R.error("无法打开打印窗口，请检查浏览器设置")};return I(()=>S.modelValue,i=>{_.value=i,i&&f()}),I(_,i=>{i||E("update:modelValue",!1)}),(i,r)=>{const u=p("el-button"),h=p("el-tag"),x=p("el-descriptions-item"),M=p("el-descriptions"),$=p("el-divider"),V=p("el-table-column"),U=p("el-table");return g(),k(ae,{modelValue:_.value,"onUpdate:modelValue":r[0]||(r[0]=j=>_.value=j),title:"入库单详情",width:"900px",loading:c.value,onCancel:A},{extra:t(()=>[e(u,{type:"primary",icon:D(he),onClick:v,size:"small"},{default:t(()=>[...r[1]||(r[1]=[s(" 打印 ",-1)])]),_:1},8,["icon"])]),default:t(()=>[e(M,{column:2,border:""},{default:t(()=>[e(x,{label:"入库单号",span:2},{default:t(()=>[e(h,{type:"primary"},{default:t(()=>[s(d(a.value.order_no||"-"),1)]),_:1})]),_:1}),e(x,{label:"状态"},{default:t(()=>[e(h,{type:N(a.value.status),size:"small"},{default:t(()=>[s(d(J(a.value.status)),1)]),_:1},8,["type"])]),_:1}),e(x,{label:"供应商"},{default:t(()=>[s(d(a.value.supplier||"-"),1)]),_:1}),e(x,{label:"联系人"},{default:t(()=>[s(d(a.value.contact||"-"),1)]),_:1}),e(x,{label:"验收人"},{default:t(()=>[s(d(a.value.receiver||"-"),1)]),_:1}),e(x,{label:"创建人"},{default:t(()=>[s(d(a.value.creator_name||"-"),1)]),_:1}),e(x,{label:"创建时间"},{default:t(()=>[s(d(a.value.created_at||"-"),1)]),_:1}),e(x,{label:"关联项目",span:2},{default:t(()=>[s(d(a.value.project_name||"-"),1)]),_:1}),a.value.remark?(g(),k(x,{key:0,label:"备注",span:2},{default:t(()=>[s(d(a.value.remark||"-"),1)]),_:1})):q("",!0)]),_:1}),G("div",dt,[a.value.status?(g(),k(we,{key:0,status:a.value.status,"status-time":a.value.updated_at||a.value.created_at,"status-description":F(a.value.status),"workflow-type":"inbound"},null,8,["status","status-time","status-description"])):q("",!0)]),e($,{"content-position":"left"},{default:t(()=>[...r[2]||(r[2]=[s("物资明细",-1)])]),_:1}),e(U,{data:a.value.items,border:"",stripe:"",style:{width:"100%"},size:"small"},{default:t(()=>[e(V,{prop:"material_name",label:"物资名称","min-width":"150","show-overflow-tooltip":""}),e(V,{prop:"spec",label:"规格型号",width:"120","show-overflow-tooltip":""}),e(V,{prop:"material",label:"材质",width:"100","show-overflow-tooltip":""}),e(V,{prop:"unit",label:"单位",width:"80"}),e(V,{prop:"quantity",label:"数量",width:"100",align:"right"}),e(V,{prop:"unit_price",label:"单价",width:"100",align:"right"},{default:t(j=>[s(d(j.row.unit_price?L(j.row.unit_price):"-"),1)]),_:1}),e(V,{label:"金额",width:"120",align:"right"},{default:t(j=>[s(d(L((j.row.quantity||0)*(j.row.unit_price||0))),1)]),_:1}),e(V,{prop:"remark",label:"备注","min-width":"120","show-overflow-tooltip":""})]),_:1},8,["data"]),T.value.length>0?(g(),Y(te,{key:0},[e($,{"content-position":"left"},{default:t(()=>[...r[3]||(r[3]=[s("审批历史",-1)])]),_:1}),e(xe,{histories:T.value},null,8,["histories"])],64)):q("",!0)]),_:1},8,["modelValue","loading"])}}},pt=de(ut,[["__scopeId","data-v-0fd3f22d"]]),ct={key:1},ft={style:{"margin-top":"20px"}},mt={__name:"RequisitionDetailDialog",props:{modelValue:{type:Boolean,default:!1},requisitionNo:{type:String,default:""}},emits:["update:modelValue"],setup(X,{emit:H}){const S=X,E=H,_=b(!1),c=b(!1),a=b({requisition_no:"",applicant:"",applicant_name:"",department:"",project_id:null,project_name:"",requisition_date:"",purpose:"",urgent:!1,remark:"",items:[],items_count:0,status:"",approved_by:"",approved_at:"",issued_by:"",issued_at:"",created_at:"",updated_at:""}),T=b([]),L=v=>({pending:"info",approved:"warning",rejected:"danger",issued:"success"})[v]||"info",N=v=>({pending:"待审核",approved:"已批准",rejected:"已拒绝",issued:"已发放"})[v]||v,J=v=>({pending:"出库单等待审核",approved:"出库单已批准，等待发放",rejected:"出库单已被拒绝",issued:"出库单已发放完成"})[v]||"",F=async v=>{try{const i=await re.getWorkflowHistory(v);T.value=i.data||[]}catch(i){console.error("获取工作流历史失败:",i)}},P=async()=>{if(S.requisitionNo){c.value=!0;try{const i=(await re.getList({pageSize:1e3})).data.find(u=>u.requisition_no===S.requisitionNo);if(!i){R.error("未找到该出库单"),f();return}const r=await re.getDetail(i.id);a.value=r.data||r,await F(i.id)}catch(v){console.error("加载出库单详情失败:",v),R.error("加载出库单详情失败"),f()}finally{c.value=!1}}},f=()=>{_.value=!1,E("update:modelValue",!1)},A=()=>{var r;const v=`
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="utf-8">
      <title>出库单 - ${a.value.requisition_no}</title>
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
          <span class="info-value">${a.value.requisition_no||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">状态：</span>
          <span class="info-value">${N(a.value.status)}</span>
        </div>
        <div class="info-row">
          <span class="info-label">项目名称：</span>
          <span class="info-value">${a.value.project_name||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">申请人：</span>
          <span class="info-value">${a.value.applicant_name||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">部门：</span>
          <span class="info-value">${a.value.department||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">申请日期：</span>
          <span class="info-value">${a.value.requisition_date||"-"}</span>
        </div>
        <div class="info-row">
          <span class="info-label">用途：</span>
          <span class="info-value">${a.value.purpose||"-"}</span>
        </div>
        ${a.value.urgent?`
        <div class="info-row">
          <span class="info-label">紧急：</span>
          <span class="info-value">是</span>
        </div>
        `:""}
        ${a.value.approved_by?`
        <div class="info-row">
          <span class="info-label">审核人：</span>
          <span class="info-value">${a.value.approved_by}</span>
        </div>
        `:""}
        ${a.value.approved_at?`
        <div class="info-row">
          <span class="info-label">审核时间：</span>
          <span class="info-value">${a.value.approved_at}</span>
        </div>
        `:""}
        ${a.value.issued_by?`
        <div class="info-row">
          <span class="info-label">发货人：</span>
          <span class="info-value">${a.value.issued_by}</span>
        </div>
        `:""}
        ${a.value.issued_at?`
        <div class="info-row">
          <span class="info-label">发货时间：</span>
          <span class="info-value">${a.value.issued_at}</span>
        </div>
        `:""}
        ${a.value.remark?`
        <div class="info-row">
          <span class="info-label">备注：</span>
          <span class="info-value">${a.value.remark}</span>
        </div>
        `:""}
      </div>

      <div class="section-title">物资明细 (${a.value.items_count||((r=a.value.items)==null?void 0:r.length)||0})</div>
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
          ${a.value.items.map(u=>`
            <tr>
              <td>${u.material||"-"}</td>
              <td class="text-left">${u.material_name||"-"}</td>
              <td>${u.specification||"-"}</td>
              <td>${u.unit||"-"}</td>
              <td>${u.requested_quantity||0}</td>
              <td>${u.approved_quantity||0}</td>
            </tr>
          `).join("")}
        </tbody>
      </table>

      <div style="margin-top: 30px; text-align: right; font-size: 10px; color: #999;">
        打印时间：${new Date().toLocaleString("zh-CN")}
      </div>
    </body>
    </html>
  `,i=window.open("","_blank");i?(i.document.write(v),i.document.close(),i.onload=()=>{i.print(),i.close()}):R.error("无法打开打印窗口，请检查浏览器设置")};return I(()=>S.modelValue,v=>{_.value=v,v&&P()}),I(_,v=>{v||E("update:modelValue",!1)}),(v,i)=>{const r=p("el-button"),u=p("el-tag"),h=p("el-descriptions-item"),x=p("el-descriptions"),M=p("el-divider"),$=p("el-table-column"),V=p("el-table");return g(),k(ae,{modelValue:_.value,"onUpdate:modelValue":i[0]||(i[0]=U=>_.value=U),title:"出库单详情",width:"900px",loading:c.value,onCancel:f},{extra:t(()=>[e(r,{type:"primary",icon:D(he),onClick:A,size:"small"},{default:t(()=>[...i[1]||(i[1]=[s(" 打印 ",-1)])]),_:1},8,["icon"])]),default:t(()=>[e(x,{column:2,border:""},{default:t(()=>[e(h,{label:"出库单号",span:2},{default:t(()=>[e(u,{type:"primary"},{default:t(()=>[s(d(a.value.requisition_no||"-"),1)]),_:1})]),_:1}),e(h,{label:"状态"},{default:t(()=>[e(u,{type:L(a.value.status),size:"small"},{default:t(()=>[s(d(N(a.value.status)),1)]),_:1},8,["type"])]),_:1}),e(h,{label:"紧急"},{default:t(()=>[a.value.urgent?(g(),k(u,{key:0,type:"danger",size:"small"},{default:t(()=>[...i[2]||(i[2]=[s("紧急",-1)])]),_:1})):(g(),Y("span",ct,"否"))]),_:1}),e(h,{label:"项目名称",span:2},{default:t(()=>[s(d(a.value.project_name||"-"),1)]),_:1}),e(h,{label:"申请人"},{default:t(()=>[s(d(a.value.applicant_name||"-"),1)]),_:1}),e(h,{label:"部门"},{default:t(()=>[s(d(a.value.department||"-"),1)]),_:1}),e(h,{label:"申请日期"},{default:t(()=>[s(d(a.value.requisition_date||"-"),1)]),_:1}),e(h,{label:"用途"},{default:t(()=>[s(d(a.value.purpose||"-"),1)]),_:1}),e(h,{label:"创建时间"},{default:t(()=>[s(d(a.value.created_at||"-"),1)]),_:1}),a.value.approved_by?(g(),k(h,{key:0,label:"审核人"},{default:t(()=>[s(d(a.value.approved_by||"-"),1)]),_:1})):q("",!0),a.value.approved_at?(g(),k(h,{key:1,label:"审核时间"},{default:t(()=>[s(d(a.value.approved_at||"-"),1)]),_:1})):q("",!0),a.value.issued_by?(g(),k(h,{key:2,label:"发货人"},{default:t(()=>[s(d(a.value.issued_by||"-"),1)]),_:1})):q("",!0),a.value.issued_at?(g(),k(h,{key:3,label:"发货时间"},{default:t(()=>[s(d(a.value.issued_at||"-"),1)]),_:1})):q("",!0),e(h,{label:"备注",span:2},{default:t(()=>[s(d(a.value.remark||"-"),1)]),_:1})]),_:1}),G("div",ft,[a.value.status?(g(),k(we,{key:0,status:a.value.status,"status-time":a.value.updated_at||a.value.created_at,"status-description":J(a.value.status),"workflow-type":"requisition"},null,8,["status","status-time","status-description"])):q("",!0)]),e(M,{"content-position":"left"},{default:t(()=>{var U;return[s("物资明细 ("+d(a.value.items_count||((U=a.value.items)==null?void 0:U.length)||0)+")",1)]}),_:1}),e(V,{data:a.value.items,border:"",stripe:"",style:{width:"100%"},size:"small"},{default:t(()=>[e($,{prop:"material",label:"材质",width:"100","show-overflow-tooltip":""}),e($,{prop:"material_name",label:"物资名称","min-width":"150","show-overflow-tooltip":""}),e($,{prop:"specification",label:"规格型号","min-width":"150","show-overflow-tooltip":""}),e($,{prop:"unit",label:"单位",width:"80"}),e($,{prop:"requested_quantity",label:"申请数量",width:"100",align:"right"}),e($,{prop:"approved_quantity",label:"批准数量",width:"100",align:"right"})]),_:1},8,["data"]),T.value.length>0?(g(),Y(te,{key:0},[e(M,{"content-position":"left"},{default:t(()=>[...i[3]||(i[3]=[s("审批历史",-1)])]),_:1}),e(xe,{histories:T.value},null,8,["histories"])],64)):q("",!0)]),_:1},8,["modelValue","loading"])}}},vt=de(mt,[["__scopeId","data-v-5e654ede"]]),_t={class:"stock-container"},gt={style:{float:"left"}},bt={style:{float:"right",color:"#8492a6","font-size":"13px"}},yt={key:1},ht={__name:"Stock",setup(X){const H=Qe(),S=b(!1),E=b([]),_=ee({page:1,pageSize:20,total:0}),c=ee({keyword:"",category:"",project_id:"",status:""}),a=b([]),T=b([]),L=b(!1),N=b("in"),J=ge(()=>({in:"库存入库",out:"库存出库",adjust:"库存调整"})[N.value]||"库存操作"),F=b(!1),P=b(null),f=ee({material_id:null,quantity:null,price:null,type:"in",remark:""}),A=b([]),v=ge(()=>A.value.find(l=>l.id===f.material_id)),i={material_id:[{required:!0,message:"请选择物资",trigger:"change"}],quantity:[{required:!0,message:"请输入数量",trigger:"blur"}],type:[{required:!0,message:"请选择操作类型",trigger:"change"}]},r=b(!1),u=b(!1),h=b([]),x=ee({page:1,pageSize:20,total:0}),M=b(null),$=b(!1),V=b(!1),U=b(""),j=b(""),W=async()=>{S.value=!0;try{let l=[];c.project_id&&(l=ke(c.project_id,a.value));const o={page:_.page,page_size:_.pageSize,search:c.keyword||void 0,category:c.category||void 0,project_ids:l.length>0?l.join(","):void 0,status:c.status||void 0},{data:m,pagination:y}=await Q.getList(o);E.value=m||[],_.total=(y==null?void 0:y.total)||0}catch(l){console.error("获取库存列表失败:",l)}finally{S.value=!1}},ke=(l,o)=>{const m=[l],y=z=>{for(const C of z){if(C.id===l){const O=w=>{if(w.children&&w.children.length>0)for(const B of w.children)m.push(B.id),O(B)};return O(C),!0}if(C.children&&C.children.length>0&&y(C.children))return!0}return!1};return y(o),m},Ve=async()=>{try{const{projectApi:l}=await _e(async()=>{const{projectApi:m}=await import("./index-BAziDdcB.js").then(y=>y.l);return{projectApi:m}},__vite__mapDeps([0,1,2,3])),{data:o}=await l.getList({pageSize:1e3});a.value=$e(o||[])}catch(l){console.error("获取项目列表失败:",l)}},$e=l=>{if(!l||l.length===0)return[];const o=new Map;l.forEach(y=>{o.set(y.id,{...y,children:[]})});const m=[];return l.forEach(y=>{const z=o.get(y.id);if(!y.parent_id)m.push(z);else{const C=o.get(y.parent_id);C?C.children.push(z):m.push(z)}}),m},ze=async()=>{try{const{materialApi:l}=await _e(async()=>{const{materialApi:m}=await import("./index-BAziDdcB.js").then(y=>y.l);return{materialApi:m}},__vite__mapDeps([0,1,2,3])),{data:o}=await l.getCategories();T.value=o||[]}catch(l){console.error("获取物资分类列表失败:",l)}},qe=()=>{_.page=1,W()},Se=()=>{c.keyword="",c.category="",c.project_id="",c.status="",_.page=1,W()},Ce=()=>{N.value="in",pe(),f.type="in",L.value=!0,ue()},je=()=>{N.value="out",pe(),f.type="out",L.value=!0,ue()},De=()=>{const l=v.value;l&&(f.price=l.price)},Te=async()=>{if(P.value)try{await P.value.validate(),F.value=!0;const l={material_id:f.material_id,quantity:f.quantity,price:f.price,type:f.type,remark:f.remark};f.type==="in"||f.type==="adjust"?(await Q.in(l),R.success("入库成功")):(await Q.out(l),R.success("出库成功")),L.value=!1,W()}catch(l){console.error("提交失败:",l)}finally{F.value=!1}},Le=l=>{M.value=l.id,r.value=!0,le()},ue=async()=>{try{const{data:l}=await Ge.getList({pageSize:1e3});A.value=l||[]}catch(l){console.error("获取物资列表失败:",l)}},le=async()=>{if(M.value){u.value=!0;try{const l={page:x.page,page_size:x.pageSize,stock_id:M.value},{data:o,pagination:m}=await Q.getLogs(l);h.value=o||[],x.total=(m==null?void 0:m.total)||0}catch(l){console.error("获取库存日志失败:",l)}finally{u.value=!1}}},Ue=async()=>{try{const l=await Q.export(c),o=new Blob([l],{type:"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}),m=window.URL.createObjectURL(o),y=document.createElement("a");y.href=m,y.download=`库存列表_${new Date().getTime()}.xlsx`,y.click(),window.URL.revokeObjectURL(m),R.success("导出成功")}catch(l){console.error("导出失败:",l)}},pe=()=>{Object.assign(f,{material_id:null,quantity:null,price:null,type:"in",remark:""}),P.value&&P.value.clearValidate()},Ne=l=>Number(l).toLocaleString("zh-CN",{minimumFractionDigits:2,maximumFractionDigits:2}),Pe=l=>l<=0?"danger":l<10?"warning":"success",Ae=l=>{const o=l.quantity||0,m=l.safe_stock||0;return o<=0?"danger":o<m?"warning":"success"},Re=l=>{const o=l.quantity||0,m=l.safe_stock||0;return o<=0?"库存不足":o<m?"库存偏低":"正常"},Ee=l=>({in:"success",out:"warning",adjust:"info"})[l]||"info",Me=l=>({in:"入库",out:"出库",adjust:"调整"})[l]||l,Oe=l=>l?l.includes("入库单")||l.includes("出库单发放"):!1,Ie=l=>{if(l.inbound_code){U.value=l.inbound_code,$.value=!0;return}if(l.requisition_code){j.value=l.requisition_code,V.value=!0;return}if(!l.remark)return;const o=l.remark.match(/入库单审核入库[:：]\s*(\w+)/);if(o){U.value=o[1],$.value=!0;return}const m=l.remark.match(/出库单发放[:：]\s*(\w+)/);if(m){j.value=m[1],V.value=!0;return}};let oe=null;const Z=()=>{oe&&clearTimeout(oe),oe=setTimeout(()=>{_.page=1,W()},500)},He=l=>{_.page=l,W()},Fe=l=>{_.pageSize=l,_.page=1,W()};return I(()=>c.keyword,Z),I(()=>c.category,Z),I(()=>c.project_id,Z),I(()=>c.status,Z),Xe(()=>{Ve(),ze(),W()}),(l,o)=>{const m=p("el-icon"),y=p("el-input"),z=p("el-option"),C=p("el-select"),O=p("el-button"),w=p("el-table-column"),B=p("el-tag"),ce=p("el-table"),fe=p("el-pagination"),We=p("el-card"),K=p("el-form-item"),Be=p("el-alert"),me=p("el-input-number"),ne=p("el-radio"),Ye=p("el-radio-group"),Je=p("el-form"),Ke=p("el-link"),ve=Ze("loading");return g(),Y("div",_t,[e(We,{shadow:"never"},{default:t(()=>[e(st,null,{left:t(()=>[e(rt,{modelValue:c.project_id,"onUpdate:modelValue":o[0]||(o[0]=n=>c.project_id=n),projects:a.value,placeholder:"选择项目（支持层级显示）",width:"300px"},null,8,["modelValue","projects"]),e(y,{modelValue:c.keyword,"onUpdate:modelValue":o[1]||(o[1]=n=>c.keyword=n),placeholder:"搜索物资名称、编码",clearable:"",style:{width:"250px"},onKeyup:et(qe,["enter"])},{prefix:t(()=>[e(m,null,{default:t(()=>[e(D(ot))]),_:1})]),_:1},8,["modelValue"]),e(C,{modelValue:c.category,"onUpdate:modelValue":o[2]||(o[2]=n=>c.category=n),placeholder:"物资分类",clearable:"",style:{width:"150px"}},{default:t(()=>[e(z,{label:"全部",value:""}),(g(!0),Y(te,null,be(T.value,n=>(g(),k(z,{key:n.id,label:n.name,value:n.name},null,8,["label","value"]))),128))]),_:1},8,["modelValue"]),e(C,{modelValue:c.status,"onUpdate:modelValue":o[3]||(o[3]=n=>c.status=n),placeholder:"库存状态",clearable:"",style:{width:"150px"}},{default:t(()=>[e(z,{label:"所有状态",value:""}),e(z,{label:"正常",value:"normal"}),e(z,{label:"库存偏低",value:"low"}),e(z,{label:"库存不足",value:"shortage"})]),_:1},8,["modelValue"]),e(O,{icon:D(nt),onClick:Se},{default:t(()=>[...o[17]||(o[17]=[s("重置",-1)])]),_:1},8,["icon"])]),right:t(()=>[D(H).hasPermission("stock_in")?(g(),k(O,{key:0,type:"primary",icon:D(tt),onClick:Ce},{default:t(()=>[...o[18]||(o[18]=[s(" 入库 ",-1)])]),_:1},8,["icon"])):q("",!0),D(H).hasPermission("stock_out")?(g(),k(O,{key:1,type:"warning",icon:D(at),onClick:je},{default:t(()=>[...o[19]||(o[19]=[s(" 出库 ",-1)])]),_:1},8,["icon"])):q("",!0),D(H).hasPermission("stock_export")?(g(),k(O,{key:2,type:"success",icon:D(lt),onClick:Ue},{default:t(()=>[...o[20]||(o[20]=[s(" 导出 ",-1)])]),_:1},8,["icon"])):q("",!0)]),_:1}),ye((g(),k(ce,{data:E.value,border:"",stripe:"",style:{width:"100%"}},{default:t(()=>[e(w,{prop:"material_code",label:"物资编码",width:"130"}),e(w,{prop:"material_name",label:"物资名称","min-width":"150","show-overflow-tooltip":""}),e(w,{prop:"category",label:"分类",width:"100"},{default:t(n=>[e(B,{size:"small"},{default:t(()=>[s(d(n.row.category||"-"),1)]),_:2},1024)]),_:1}),e(w,{prop:"specification",label:"规格型号",width:"120","show-overflow-tooltip":""}),e(w,{prop:"unit",label:"单位",width:"80"}),e(w,{prop:"quantity",label:"库存数量",width:"120",align:"right"},{default:t(n=>[e(B,{type:Pe(n.row.quantity),size:"large"},{default:t(()=>[s(d(n.row.quantity||0),1)]),_:2},1032,["type"])]),_:1}),e(w,{prop:"safety_stock",label:"安全库存",width:"100",align:"right"},{default:t(n=>[s(d(n.row.safety_stock||"-"),1)]),_:1}),e(w,{prop:"stock_status",label:"库存状态",width:"100"},{default:t(n=>[e(B,{type:Ae(n.row),size:"small"},{default:t(()=>[s(d(Re(n.row)),1)]),_:2},1032,["type"])]),_:1}),e(w,{prop:"project_name",label:"关联项目","min-width":"150","show-overflow-tooltip":""},{default:t(n=>[s(d(n.row.project_name||"-"),1)]),_:1}),e(w,{prop:"updated_at",label:"最后更新",width:"160"},{default:t(n=>[s(d(n.row.updated_at||"-"),1)]),_:1}),e(w,{label:"操作",width:"150",fixed:"right"},{default:t(n=>[e(O,{type:"primary",size:"small",icon:D(it),onClick:ie=>Le(n.row)},{default:t(()=>[...o[21]||(o[21]=[s(" 日志 ",-1)])]),_:1},8,["icon","onClick"])]),_:1})]),_:1},8,["data"])),[[ve,S.value]]),e(fe,{"current-page":_.page,"onUpdate:currentPage":o[4]||(o[4]=n=>_.page=n),"page-size":_.pageSize,"onUpdate:pageSize":o[5]||(o[5]=n=>_.pageSize=n),"page-sizes":[10,20,50,100],total:_.total,layout:"total, sizes, prev, pager, next, jumper",onSizeChange:Fe,onCurrentChange:He,class:"mt-20"},null,8,["current-page","page-size","total"])]),_:1}),e(ae,{modelValue:L.value,"onUpdate:modelValue":o[11]||(o[11]=n=>L.value=n),title:J.value,width:"700px",loading:F.value,onConfirm:Te},{default:t(()=>[e(Je,{ref_key:"formRef",ref:P,model:f,rules:i,"label-width":"100px"},{default:t(()=>[e(K,{label:"物资",prop:"material_id"},{default:t(()=>[e(C,{modelValue:f.material_id,"onUpdate:modelValue":o[6]||(o[6]=n=>f.material_id=n),placeholder:"请选择物资",filterable:"",style:{width:"100%"},onChange:De},{default:t(()=>[(g(!0),Y(te,null,be(A.value,n=>(g(),k(z,{key:n.id,label:`${n.code} - ${n.name}`,value:n.id},{default:t(()=>[G("span",gt,d(n.code)+" - "+d(n.name),1),G("span",bt,d(n.category),1)]),_:2},1032,["label","value"]))),128))]),_:1},8,["modelValue"])]),_:1}),v.value?(g(),k(Be,{key:0,title:`当前库存: ${v.value.quantity} ${v.value.unit}`,type:"info",closable:!1,class:"mb-2"},null,8,["title"])):q("",!0),e(K,{label:"数量",prop:"quantity"},{default:t(()=>{var n;return[e(me,{modelValue:f.quantity,"onUpdate:modelValue":o[7]||(o[7]=ie=>f.quantity=ie),min:1,max:N.value==="out"?(n=v.value)==null?void 0:n.quantity:void 0,step:1,precision:0,placeholder:"请输入数量",style:{width:"100%"}},null,8,["modelValue","max"])]}),_:1}),e(K,{label:"单价",prop:"price"},{default:t(()=>[e(me,{modelValue:f.price,"onUpdate:modelValue":o[8]||(o[8]=n=>f.price=n),min:0,precision:2,step:.01,placeholder:"请输入单价",style:{width:"100%"}},null,8,["modelValue"]),o[22]||(o[22]=G("div",{class:"text-gray"},"不填则使用物资默认单价",-1))]),_:1}),e(K,{label:"操作类型",prop:"type"},{default:t(()=>[e(Ye,{modelValue:f.type,"onUpdate:modelValue":o[9]||(o[9]=n=>f.type=n)},{default:t(()=>[e(ne,{label:"in"},{default:t(()=>[...o[23]||(o[23]=[s("入库",-1)])]),_:1}),e(ne,{label:"out"},{default:t(()=>[...o[24]||(o[24]=[s("出库",-1)])]),_:1}),e(ne,{label:"adjust"},{default:t(()=>[...o[25]||(o[25]=[s("调整",-1)])]),_:1})]),_:1},8,["modelValue"])]),_:1}),e(K,{label:"备注",prop:"remark"},{default:t(()=>[e(y,{modelValue:f.remark,"onUpdate:modelValue":o[10]||(o[10]=n=>f.remark=n),type:"textarea",rows:3,placeholder:"请输入备注",maxlength:"500"},null,8,["modelValue"])]),_:1})]),_:1},8,["model"])]),_:1},8,["modelValue","title","loading"]),e(ae,{modelValue:r.value,"onUpdate:modelValue":o[14]||(o[14]=n=>r.value=n),title:"库存日志",width:"900px","show-footer":!1},{default:t(()=>[ye((g(),k(ce,{data:h.value,border:"",stripe:"","max-height":"400"},{default:t(()=>[e(w,{prop:"created_at",label:"时间",width:"160"}),e(w,{prop:"type",label:"类型",width:"80"},{default:t(n=>[e(B,{type:Ee(n.row.type),size:"small"},{default:t(()=>[s(d(Me(n.row.type)),1)]),_:2},1032,["type"])]),_:1}),e(w,{prop:"quantity",label:"数量",width:"80",align:"right"}),e(w,{prop:"quantity_before",label:"操作前",width:"80",align:"right"}),e(w,{prop:"quantity_after",label:"操作后",width:"80",align:"right"}),e(w,{prop:"price",label:"单价",width:"80",align:"right"},{default:t(n=>[s(d(n.row.price?Ne(n.row.price):"-"),1)]),_:1}),e(w,{prop:"remark",label:"备注","min-width":"200"},{default:t(n=>[n.row.inbound_code||n.row.requisition_code||Oe(n.row.remark)?(g(),k(Ke,{key:0,type:"primary",onClick:ie=>Ie(n.row)},{default:t(()=>[s(d(n.row.remark),1)]),_:2},1032,["onClick"])):(g(),Y("span",yt,d(n.row.remark||"-"),1))]),_:1})]),_:1},8,["data"])),[[ve,u.value]]),e(fe,{"current-page":x.page,"onUpdate:currentPage":o[12]||(o[12]=n=>x.page=n),"page-size":x.pageSize,"onUpdate:pageSize":o[13]||(o[13]=n=>x.pageSize=n),"page-sizes":[10,20,50],total:x.total,layout:"total, sizes, prev, pager, next",onSizeChange:le,onCurrentChange:le,class:"mt-20"},null,8,["current-page","page-size","total"])]),_:1},8,["modelValue"]),e(pt,{modelValue:$.value,"onUpdate:modelValue":o[15]||(o[15]=n=>$.value=n),"order-no":U.value},null,8,["modelValue","order-no"]),e(vt,{modelValue:V.value,"onUpdate:modelValue":o[16]||(o[16]=n=>V.value=n),"requisition-no":j.value},null,8,["modelValue","requisition-no"])])}}},Ct=de(ht,[["__scopeId","data-v-f061febc"]]);export{Ct as default};
