select e.id, e.amount, e.product, e.shop, e.expend_date, c.id, c.name, eu.user_id
from expenses e
join categories c on e.category_id = c.id
left join expense_users eu on e.id = eu.expense_id
where ()
order by e.expend_date desc
limit 35
offset 0

