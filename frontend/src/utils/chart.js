export default {
    randomColor() {
        return `#${Math.floor(Math.random() * 0xffffff).toString(16).padEnd(6, "0")}`;
    },
    getYAxis(datas) {
        const getValuePow = (value) => {
            let pow = 0;
            while(value >= 40) {
                value /= 10;
                pow++;
            }
            while(value < 4 && value > 0) {
                value *= 10;
                pow--;
            }
            return {value, pow};
        };
        
        const max_value = Math.max(...datas);
        const min_value = Math.min(...datas);
        let temp_value = 0;
        if(max_value > 0 && min_value > 0) {
            temp_value = max_value;
        } else if(max_value < 0 && min_value < 0) {
            temp_value = -min_value;
        } else if(max_value == 0 && min_value == 0) {
            temp_value = 1;
        } else {
            temp_value = max_value - min_value;
        }
        
        const {value, pow} = getValuePow(temp_value);
        
        const steps = [1, 2, 3, 5];
        const lines = [4, 5, 6, 7, 8];
        const except = [14, 16, 21, 24];
        let step = 1;
        let line = 5;
        loop: for(let s of steps) {
            for(let l of lines) {
                if(s * l >= value) {
                    step = s;
                    line = l;
                    if(!~except.indexOf(s * l)) {
                        break loop;
                    }
                }
            }
        }

        const list = [];
        const step_value = Math.pow(10, pow) * step;
        //偏移
        let offset = 0;
        while(min_value < offset) {
            offset -= step_value;
        }

        while(line > 0) {
            let axis = step_value * line;
            if(pow < 0) {
                axis = parseFloat(axis.toFixed(-pow));
            }
            
            list.push(axis + offset);
            line--;
        }
        list.push(0 + offset);

        //修正
        while(list[0] < max_value) {
            list.unshift(list[0] + step_value);
        }
        return list;
    },
    getFontSize(width, dpr, normal = 8, max = 12) {
        let font_size = width / 132 * normal;
        if(font_size > max * dpr) {
            font_size = max * dpr;
        }
        return font_size;
    },
    drawYAxis(canvas, y_axis) {
        const total = y_axis.length - 1;
        const width = canvas.width;
        const height = canvas.height;
        
        const ctx = canvas.getContext("2d");
        ctx.save();
        ctx.strokeStyle = '#eee';
        ctx.lineWidth = 1;
        y_axis.forEach((v, i) => {
            ctx.beginPath();
            const y = parseInt(i * height / total) + (i == 0 ? 0.5 : -0.5);
            ctx.lineTo(0, y);
            ctx.lineTo(width, y);
            ctx.stroke();
        });
        ctx.restore();
    },
    drawYLabels(canvas, y_axis) {
        const dpr = canvas.dpr;
        const total = y_axis.length - 1;
        const width = canvas.width;
        const height = canvas.height;
        const font_size = this.getFontSize(width, dpr, 6, 10);
        
        const ctx = canvas.getContext("2d");
        ctx.save();
        ctx.globalAlpha = 0.5;
        ctx.font = font_size + "px Arial";
        ctx.fillStyle = "#666";
        ctx.textAlign = 'start';
        ctx.strokeStyle = "#fff";
        ctx.lineWidth = 2 * dpr;
        y_axis.forEach((v, i) => {
            const y = parseInt(height - (total - i) * height / total) + (i == 0 ? dpr : 0);
            ctx.textBaseline = i == 0 ? 'top' : i == total ? 'bottom' : 'middle';
            ctx.strokeText(v, 0, y);
            ctx.fillText(v, 0, y);
        });
        ctx.restore();
    },
    drawXLabels(canvas, lables) {
        const dpr = canvas.dpr;
        const width = canvas.width;
        const height = canvas.height;
        const x_width = width / lables.length;
        const font_size = this.getFontSize(width, dpr, 6);

        const ctx = canvas.getContext("2d");
        ctx.save();
        ctx.globalAlpha = 0.8;
        ctx.font = font_size + "px Arial";
        ctx.strokeStyle = "#fff";
        ctx.fillStyle = "#333";
        ctx.textAlign = 'center';
        ctx.textBaseline = 'bottom';
        ctx.lineWidth = 2 * dpr;
        lables.forEach((v, i) => {
            const x = i * x_width + x_width / 2;
            const y = height - 2 * dpr;
            ctx.strokeText(v, x, y);
            ctx.fillText(v, x, y);
        });
        ctx.restore();
    },
    drawBarChart(canvas, y_axis, datas) {
        const max_axis = Math.max(...y_axis);
        const min_axis = Math.min(...y_axis);
        const dpr = canvas.dpr;
        const width = canvas.width;
        const height = canvas.height;
        const x_width = width / datas.length;
        let bar_width = x_width * 0.8;
        if(bar_width > 30 * dpr) {
            bar_width = 30 * dpr;
        }
        const left = (x_width - bar_width) / 2;
        const font_size = this.getFontSize(width, dpr);
        const ctx = canvas.getContext("2d");
        ctx.textAlign = 'center';
        ctx.strokeStyle = "#fff";
        ctx.lineWidth = 2 * dpr;
        datas.forEach((value, i) => {
            const x = i * x_width + left;
            const y = height - (value - min_axis) * height / (max_axis - min_axis);
            const bar_height = height * value / (max_axis - min_axis);
            ctx.fillStyle = i == 0 ? "#46bc99" : this.randomColor();
            ctx.fillRect(x, y, bar_width, bar_height);

            if(value == 0) {
                return;
            }

            if(value != 0) {
                ctx.save();
                const x_text = i * x_width + x_width / 2;
                const y_text = value >= 0 ? y - 1 * dpr : y + 1 * dpr;
                ctx.translate(x_text, y_text);
                ctx.font = font_size + "px Arial";
                ctx.fillStyle = "#333";
                ctx.textBaseline = value >= 0 ? 'bottom' : 'top';
                ctx.strokeText(value, 0, 0);
                ctx.fillText(value, 0, 0);
                ctx.restore();
            }
        });
    },
    drawLineChart(canvas, y_axis, datas) {
        const max_axis = Math.max(...y_axis);
        const min_axis = Math.min(...y_axis);
        const dpr = canvas.dpr;
        const width = canvas.width;
        const height = canvas.height;
        const x_width = width / datas.length;
        const radius = 3 * dpr;
        const font_size = this.getFontSize(width, dpr);
        const ctx = canvas.getContext("2d");
        ctx.beginPath();
        ctx.fillStyle = "#46bc99";
        ctx.strokeStyle = '#46bc99';
        datas.forEach((value, i) => {
            const x = i * x_width + x_width / 2;
            const y = height - (value - min_axis) * height / (max_axis - min_axis);
            ctx.lineWidth = 2 * dpr;
            ctx.lineTo(x, y);
        });
        ctx.stroke();

        datas.forEach((value, i) => {
            const x = i * x_width + x_width / 2;
            const y = height - (value - min_axis) * height / (max_axis - min_axis);
            ctx.lineWidth = 2 * dpr;
            ctx.beginPath();
            ctx.arc(x, y, radius, 0, Math.PI * 2);
            ctx.fill();

            if(value != 0) {
                ctx.save();
                const x_text = x;
                const left_value = i == 0 ? value : datas[i - 1];
                const right_value = i == datas.length - 1 ? value : datas[i + 1];
                const middle_value = (left_value + right_value) / 2;
                let y_text = value < middle_value ? y + font_size * 0.8 : y - font_size * 0.8;
                ctx.translate(x_text, y_text);
                ctx.font = font_size + "px Arial";
                ctx.fillStyle = "#333";
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                ctx.strokeStyle = "#fff";
                ctx.lineWidth = 2 * dpr;
                ctx.strokeText(value, 0, 0);
                ctx.fillText(value, 0, 0);
                ctx.restore();
            }
        });
    },
    drawPieChart(canvas, datas, labels = []) {
        const dpr = canvas.dpr;
        const width = canvas.width;
        const height = canvas.height;
        const x = width / 2;
        const y = height / 2;
        const radius = height / 2 * 0.9;
        const font_size = this.getFontSize(width, dpr);
        const label_size = this.getFontSize(width, dpr, 6);
        const ctx = canvas.getContext("2d");
        let temp_angle = -90;
        const total = datas.reduce((a, b) => Math.abs(a) + Math.abs(b));
        datas.forEach((value, i) => {
            ctx.fillStyle = total > 0 ? (i == 0 ? "#46bc99" : this.randomColor()) : '#f9f9f9';
            ctx.beginPath();
            ctx.moveTo(x, y);
            let angle = total > 0 ? Math.abs(value) / total * 360 : 360 / datas.length;
            const start_angle = temp_angle * Math.PI / 180;
            const end_angle = (temp_angle + angle) * Math.PI / 180;
            ctx.arc(x, y, radius, start_angle, end_angle);
            ctx.fill();

            //绘制文字
            var txt_angle = temp_angle + 1 / 2 * angle;
            const x_text = x + Math.cos(txt_angle * Math.PI / 180) * (radius/1.5);
            const y_text = y + Math.sin(txt_angle * Math.PI / 180) * (radius/1.5);
        
            ctx.font = font_size + "px Arial";
            ctx.fillStyle = "#333";
            ctx.textBaseline = 'middle';
            // if (txt_angle > -90 && txt_angle < 90) { //设置y轴左边的文字结束位置对齐，防止文字显示不全
            //     ctx.textAlign = 'start';
            // } else {
            //     ctx.textAlign = 'end';
            // }
            ctx.textAlign = 'center';
            ctx.strokeStyle = "#fff";
            ctx.lineWidth = 2 * dpr;
            ctx.strokeText(value, x_text, y_text);
            ctx.fillText(value, x_text, y_text); //填充文字

            //绘制标签
            if(labels[i]) {
                ctx.save();
                ctx.font = label_size + "px Arial";
                ctx.fillStyle = "#333";
                ctx.strokeStyle = "#fff";
                ctx.lineWidth = 2 * dpr;
                ctx.textAlign = 'center';
                ctx.textBaseline = 'bottom';
                let sapce = 0.5 * dpr;
                let label_rotate = txt_angle * Math.PI / 180 + Math.PI / 2;
                if(15 < txt_angle && txt_angle < 165) {
                    label_rotate -= Math.PI;
                    ctx.textBaseline = 'top';
                    sapce = 1.5 * dpr;
                }
                const x_label = x + Math.cos(txt_angle * Math.PI / 180) * (radius + sapce);
                const y_label = y + Math.sin(txt_angle * Math.PI / 180) * (radius + sapce);
                ctx.translate(x_label, y_label);
                ctx.rotate(label_rotate);
                ctx.strokeText(labels[i], 0, 0);
                ctx.fillText(labels[i], 0, 0);
                ctx.restore();
            }

            temp_angle += angle;
        });
    },
    drawChart(canvas, data) {
        canvas.dpr = window.devicePixelRatio || 1;
        const datas = data.list;
        const labels = data.label;
        if(data.type == 'bar') {
            const y_axis = this.getYAxis(datas);
            this.drawYAxis(canvas, y_axis);
            this.drawBarChart(canvas, y_axis, datas);
            this.drawYLabels(canvas, y_axis);
            this.drawXLabels(canvas, labels);
        } else if(data.type == 'line') {
            const y_axis = this.getYAxis(datas);
            this.drawYAxis(canvas, y_axis);
            this.drawLineChart(canvas, y_axis, datas);
            this.drawYLabels(canvas, y_axis);
            this.drawXLabels(canvas, labels);
        } else if(data.type == 'pie') {
            this.drawPieChart(canvas, datas, labels);
        }
    },
}