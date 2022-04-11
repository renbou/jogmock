import V from"svelte-preprocess";import*as o from"path";import*as g from"fs/promises";import*as i from"path";import*as b from"fs/promises";var A=(s,...t)=>c.merging(...Object.getOwnPropertyNames(s).map(e=>c.alias({alias:e,directory:i.resolve(s[e])},...t)));var x=s=>{let t=A(s,c.sass,c.pkgJson);return function(e,r,n){Promise.resolve(t(e)).then(l=>{if(l!==void 0&&l.length>0){n({file:l[0]});return}n(null)})}};function R(s,t,e){if(!s.startsWith("json:")){e(null);return}s=s.slice(5),b.readFile(i.join(i.dirname(t),`${s}.json`)).then(r=>{let n=JSON.parse(r.toString()),l=p=>p.replace(/([A-Z])/g," $1").split(" ").map(a=>a[0].toLowerCase()+a.slice(1)).join("-"),y="",v=(p,a)=>{for(let f in a){let h=`${p}${p&&"-"}${l(f)}`;typeof a[f]=="string"?y+=`$${h}: ${a[f]};
`:v(h,a[f])}};v(l(i.parse(s).name),n),e({contents:y})})}function I(s){return[s]}async function u(s,t){return g.stat(s).then(e=>t&&t(e)||!t?[s]:void 0).catch(()=>!1)}async function m(s){return u(s,t=>t.isFile())}async function L(s){return u(s,t=>t.isDirectory())}function d(...s){return async function(t){let e=[];for(let r of s.map(n=>`${t}.${n}`))await m(r)&&e.push(r);return e.length>0?e:void 0}}var P=d("js","cjs","mjs"),j=d("sass","scss","css");async function O(s){if(/\.(css|scss|sass)$/.test(s))return[s];let t=[s,`${s}/index`,`${s}/_index`];o.basename(s).startsWith("_")||t.push(`${o.dirname(s)}/_${o.basename(s)}`);let e=[];for(let r of t){let n=await j(r);n&&e.push(...n)}return e}async function E(s){let t=o.join(s,"package.json");try{if(await m(t)){let e=JSON.parse((await g.readFile(t)).toString()),r=e.main?o.join(s,e.main):void 0;if(r)return[r]}}catch(e){console.error(`Error resolving module package.json for path ${s}: ${e instanceof Error?e.message:e}`)}}async function D(s){let t=await P(o.join(s,"index"));if(t)return t}function $(...s){return async t=>{let e=[];for(let r of s)await Promise.resolve(r(t)).then(n=>{n!==void 0&&e.push(...n)});return e}}function J({alias:s,directory:t},...e){let r=$(...e);return async n=>{if(n.startsWith(s)&&n.length>s.length)return r(o.join(t,n.slice(s.length)))}}var c={identity:I,existing:u,existingFile:m,existingDir:L,extension:d,jsExtension:P,sassExtension:j,merging:$,sass:O,pkgJson:E,jsIndex:D,alias:J};function w(s){return s===void 0?process.env.NODE_ENV==="production":s==="production"}function S(s){return!w(s)}var k=S(),M=w();var F={preprocess:V({sourceMap:k,scss:{importer:[x({"@":"./src","~":"./node_modules"}),R],prependData:'@use "@styles/variables" as *;'}})},q=F;export{q as default};
