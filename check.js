const fs = require('fs');

function checkViews() {
  const views = fs.readdirSync('frontend/src/views').filter(f => f.endsWith('.vue'));
  for (const view of views) {
    const content = fs.readFileSync('frontend/src/views/' + view, 'utf8');
    if (!content.includes('<template>')) continue;
    const template = content.substring(content.indexOf('<template>') + 10, content.lastIndexOf('</template>'));
    
    let rootNodes = 0;
    let depth = 0;
    
    const lines = template.split('\n');
    for (let line of lines) {
      if (line.includes('<!--') && !line.includes('-->')) line = line.split('<!--')[0]; // simple ignore
      const divOpens = (line.match(/<div(\s|>|$)/gi) || []).length;
      const divCloses = (line.match(/<\/div>/gi) || []).length;
      
      if (depth === 0 && divOpens > 0) {
          rootNodes += divOpens;
      }
      depth += divOpens - divCloses;
    }
    
    console.log(view + ' | Root div nodes: ' + rootNodes + ' | Final div depth: ' + depth);
  }
}

checkViews();
