package main

import (
	"fmt"
	"html/template"
)

var funcMap = template.FuncMap{

	"charList": func(p Page) template.HTML {

		var list string
		for _, name := range p.Names {
			c := p.Data[name]
			s := fmt.Sprintf(`<option class="statListOpt" value="%s">%s level %d %s %s %s AC: %d</option>
			`,
				name, name, c["Level"], c["Ancestry"], c["Class"], c["Background"], c["AC"])
			list += s
		}
		return template.HTML(list)
	},
	"charDetails": func(p Page) template.HTML {

		var detail string
		for _, name := range p.Names {
			c := p.Data[name]
			s := fmt.Sprintf(`<div class="character_details">
								<h2>%s</h2>
							   <div class="attributes">`, name)

			detail += s

			for _, v := range p.Abilities {
				s = fmt.Sprintf(`<div class="hex">
									<div class="hex inner">
										<h2>%d</h2>
										<p><strong>%s</strong></p>
									</div>
								</div>  `, c[v], v)
				detail += s
			}
			s = fmt.Sprintf(`<div class="ac">
                                    <h2>%d</h2>
                                    <p><strong>AC</strong></p>
                                </div>
                            </div>
                            <div class="save">
                                <div class="save_inner">
                                    <h2>%d</h2>
                                    <p><strong>per</strong></p>
                                </div>
                            </div>
                            <div  class="saves">
                                <div class="save">
                                    <div class="save_inner">
                                        <h2>%d</h2>
                                        <p><strong>fort</strong></p>
                                    </div>
                                </div>
                                <div class="save">
                                    <div class="save_inner">
                                        <h2>%d</h2>
                                        <p><strong>ref</strong></p>
                                    </div>
                                </div>
                                <div class="save">
                                    <div class="save_inner">
                                        <h2>%d</h2>
                                        <p><strong>will</strong></p>
                                    </div>
                                </div>
                            </div>`, c["AC"], c["Perception"], c["Fortitude"], c["Reflex"], c["Will"])

			detail += s
			detail += `<table class="proficiencies">`

			j, k := 0, 5

			for i := 0; i < 4; i++ {
				detail += `<tr>`
				for _, v := range p.Proficiencies[j:k] {

					s = fmt.Sprintf(`<td>%s</td>
                                    <td>%d (%s)</td>
                                `, v, c[v], c[v+"Train"])

					detail += s
				}
				detail += `</tr>`
				j += 5
				k += 5
			}
			detail += "</table></div>"
		}
		return template.HTML(detail)
	}}
